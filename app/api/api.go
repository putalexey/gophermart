package api

import (
	"context"
	"errors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Api struct {
	Logger         *zap.SugaredLogger
	Address        string
	AccrualAddress string
	DatabaseDSN    string
}

func (a Api) Run(ctx context.Context) {
	// Create gin router
	router := gin.New()
	router.Use(ginzap.Ginzap(a.Logger.Desugar(), time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(a.Logger.Desugar(), true))
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("sess", store))

	// Set up GET endpoint
	router.GET("/", func(c *gin.Context) {
		//session := sessions.Default(c)
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	srv := &http.Server{
		Addr:     a.Address,
		Handler:  router,
		ErrorLog: zap.NewStdLog(a.Logger.Desugar()),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			a.Logger.Infof("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	a.Logger.Info("Shutting down api...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		a.Logger.Fatalf("Server forced to shutdown: %s", err)
	}
}
