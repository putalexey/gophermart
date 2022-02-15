package main

import (
	"context"
	"github.com/putalexey/gophermart/cmd/gophermart/config"
	"github.com/putalexey/gophermart/loyaltyapi"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalln(err)
	}

	logger := makeLogger()
	defer logger.Sync() // flushes buffer, if any

	logger.Infow(
		"starting gophermart API",
		"addr", cfg.Address,
		"accrual_addr", cfg.AccrualSystemAddress,
	)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	loyaltyAPIConfig := loyaltyapi.Config{
		DatabaseDSN:    cfg.DatabaseDSN,
		Address:        cfg.Address,
		AccrualAddress: cfg.AccrualSystemAddress,
		MigrationsDir:  cfg.MigrationsDir,
		SecretKey:      "some secret key",
	}
	app, err := loyaltyapi.New(logger, loyaltyAPIConfig)
	if err != nil {
		logger.Fatal(err)
	}
	go func() {
		// Launch Gin and
		// handle potential error
		time.Sleep(3 * time.Second)

		app.Run(ctx)

		wg.Done()
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//select {
	//case <-quit:
	//case <-ctx.Done():
	//}

	logger.Info("Gracefully shutting down server...")
	cancel()

	// wait for go routine with api ends shutdown process
	wg.Wait()

	logger.Info("Successfully shutdown")
}

func makeLogger() *zap.SugaredLogger {
	//baseLogger, _ := zap.NewProduction()
	baseLogger, _ := zap.NewDevelopment()
	return baseLogger.Sugar()
}
