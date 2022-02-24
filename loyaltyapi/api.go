package loyaltyapi

import (
	"context"
	"errors"
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/putalexey/gophermart/loyaltyapi/handlers"
	"github.com/putalexey/gophermart/loyaltyapi/models"
	"github.com/putalexey/gophermart/loyaltyapi/repository"
	"github.com/putalexey/gophermart/loyaltyapi/services"
	"github.com/putalexey/gophermart/loyaltyapi/workers"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	DatabaseDSN    string
	Address        string
	AccrualAddress string
	MigrationsDir  string
	SecretKey      string
}

type LoyaltyAPI struct {
	Logger         *zap.SugaredLogger
	DatabaseDSN    string
	Address        string
	AccrualAddress string
	MigrationsDir  string
	secretKey      string
	srv            *http.Server
	router         *gin.Engine
	repository     *repository.Repo
	accrualService *services.Accrual
}

func New(logger *zap.SugaredLogger, config Config) (*LoyaltyAPI, error) {
	app := &LoyaltyAPI{
		Logger:         logger,
		DatabaseDSN:    config.DatabaseDSN,
		Address:        config.Address,
		AccrualAddress: config.AccrualAddress,
		MigrationsDir:  config.MigrationsDir,
		secretKey:      config.SecretKey,
	}
	err := app.init()
	return app, err
}

func (a *LoyaltyAPI) init() error {
	var err error

	a.repository, err = repository.New(a.DatabaseDSN, a.MigrationsDir)
	if err != nil {
		return err
	}

	a.accrualService = &services.Accrual{
		Address: a.AccrualAddress,
	}

	// Create jwtMiddleware
	jwtMiddleware, err := ginjwt.New(&ginjwt.GinJWTMiddleware{
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: models.UserIdentityKey,
		SendCookie:  true,
		Authorizator: func(user interface{}, c *gin.Context) bool {
			return user != nil
		},
		PayloadFunc: func(data interface{}) ginjwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return ginjwt.MapClaims{
					models.UserIdentityKey: v.UUID,
				}
			}
			return ginjwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := ginjwt.ExtractClaims(c)
			user, err := a.repository.GetUser(c, claims[models.UserIdentityKey].(string))
			if err != nil {
				return nil
			}
			return user
		},
		Key:         []byte(a.secretKey),
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
	})
	if err != nil {
		return err
	}

	handle := handlers.New(a.Logger, jwtMiddleware)

	jwtMiddleware.Authenticator = handle.Authenticator(a.repository)

	// Create gin router
	a.router = gin.New()
	a.router.HandleMethodNotAllowed = true
	a.router.Use(ginzap.Ginzap(a.Logger.Desugar(), time.RFC3339, true))
	a.router.Use(ginzap.RecoveryWithZap(a.Logger.Desugar(), true))
	a.router.Use(gzip.Gzip(gzip.DefaultCompression))

	a.router.POST("/api/user/register", handle.Register(a.repository))
	a.router.POST("/api/user/login", jwtMiddleware.LoginHandler)

	authGroup := a.router.Group("")
	authGroup.Use(jwtMiddleware.MiddlewareFunc())
	authGroup.GET("/api/me", handle.CurrentUser())
	authGroup.POST("/api/user/orders", handle.UserCreateOrder(a.repository, a.repository))
	authGroup.GET("/api/user/orders", handle.UserGetOrders(a.repository))
	authGroup.GET("/api/user/balance", handle.GetUserBalance(a.repository))
	authGroup.POST("/api/user/balance/withdraw", handle.BalanceWithdraw(a.repository))
	authGroup.GET("/api/user/balance/withdrawals", handle.GetBalanceWithdrawals(a.repository))

	return nil
}

func (a *LoyaltyAPI) Run(ctx context.Context) {
	a.srv = &http.Server{
		Addr:     a.Address,
		Handler:  a.router,
		ErrorLog: zap.NewStdLog(a.Logger.Desugar()),
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			a.Logger.Infof("listen: %s\n", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		orderWorker := workers.New(ctx, a.Logger, a.repository, 10*time.Second, a.accrualService)
		orderWorker.Run()
		wg.Done()
	}()

	<-ctx.Done()

	a.Logger.Info("Shutting down api...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(shutdownCtx); err != nil {
		a.Logger.Fatalf("Server forced to shutdown: %s", err)
	}
	wg.Wait()
}
