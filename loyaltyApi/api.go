package loyaltyApi

import (
	"context"
	"database/sql"
	"errors"
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/putalexey/gophermart/loyaltyApi/handlers"
	"github.com/putalexey/gophermart/loyaltyApi/models"
	"github.com/putalexey/gophermart/loyaltyApi/repository"
	"github.com/putalexey/gophermart/loyaltyApi/responses"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func notImplementedFunc(c *gin.Context) {
	responses.JSONError(c, http.StatusNotImplemented, errors.New("not implemented"))
}

type LoyaltyApiConfig struct {
	DatabaseDSN    string
	Address        string
	AccrualAddress string
	MigrationsDir  string
	SecretKey      string
}

type LoyaltyApi struct {
	Logger         *zap.SugaredLogger
	DatabaseDSN    string
	Address        string
	AccrualAddress string
	MigrationsDir  string
	secretKey      string
	srv            *http.Server
	router         *gin.Engine
	db             *sql.DB
	repository     *repository.Repo
}

func New(logger *zap.SugaredLogger, config LoyaltyApiConfig) (*LoyaltyApi, error) {
	app := &LoyaltyApi{
		Logger:         logger,
		DatabaseDSN:    config.DatabaseDSN,
		Address:        config.Address,
		AccrualAddress: config.AccrualAddress,
		MigrationsDir:  config.MigrationsDir,
		secretKey:      config.SecretKey,
	}
	err := app.Init()
	return app, err
}

func (a *LoyaltyApi) connectToDB() error {
	db, err := sql.Open("pgx", a.DatabaseDSN)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(30 * time.Second)
	db.SetConnMaxLifetime(2 * time.Minute)

	a.db = db
	err = a.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (a *LoyaltyApi) Init() error {
	// Connect to DB
	err := a.connectToDB()
	if err != nil {
		return err
	}

	a.repository, err = repository.New(a.db, a.MigrationsDir)
	if err != nil {
		return err
	}

	// Create jwtMiddleware
	jwtMiddleware, err := ginjwt.New(&ginjwt.GinJWTMiddleware{
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		IdentityKey:   models.UserIdentityKey,
		Key:           []byte(a.secretKey),
		Authenticator: handlers.Authenticator(a.repository),
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
	})
	if err != nil {
		return err
	}

	// Create gin router
	a.router = gin.New()
	a.router.Use(ginzap.Ginzap(a.Logger.Desugar(), time.RFC3339, true))
	a.router.Use(ginzap.RecoveryWithZap(a.Logger.Desugar(), true))
	a.router.Use(gzip.Gzip(gzip.DefaultCompression))

	a.router.POST("/api/user/register", handlers.Register(a.Logger, jwtMiddleware, a.repository))
	a.router.POST("/api/user/login", handlers.Login(a.Logger, jwtMiddleware, a.repository))

	authGroup := a.router.Group("")
	authGroup.Use(jwtMiddleware.MiddlewareFunc())
	authGroup.POST("/api/user/orders", notImplementedFunc)
	authGroup.GET("/api/user/orders", notImplementedFunc)
	authGroup.GET("/api/user/balance", notImplementedFunc)
	authGroup.POST("/api/user/balance/withdraw", notImplementedFunc)
	authGroup.GET("/api/user/balance/withdrawals", notImplementedFunc)

	return nil
}

func (a *LoyaltyApi) Run(ctx context.Context) {
	a.srv = &http.Server{
		Addr:     a.Address,
		Handler:  a.router,
		ErrorLog: zap.NewStdLog(a.Logger.Desugar()),
	}

	go func() {
		if err := a.srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			a.Logger.Infof("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	a.Logger.Info("Shutting down api...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(shutdownCtx); err != nil {
		a.Logger.Fatalf("Server forced to shutdown: %s", err)
	}
}
