package main

import (
	"context"
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/putalexey/gophermart/loyaltyApi"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
)

var addressPattern = `^.+?:\d{1,5}$`

type AppConfig struct {
	Address              string `env:"RUN_ADDRESS"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseDSN          string `env:"DATABASE_URI"`
	MigrationsDir        string `env:"DATABASE_MIGRATIONS"`
}

func main() {
	cfg := AppConfig{
		Address:              "localhost:8000",
		AccrualSystemAddress: "",
		DatabaseDSN:          "",
		MigrationsDir:        "migrations",
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	parseFlags(&cfg)
	err = checkConfig(cfg)
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

	app, err := loyaltyApi.New(logger, cfg.DatabaseDSN, cfg.Address, cfg.AccrualSystemAddress, cfg.MigrationsDir)
	if err != nil {
		logger.Fatal(err)
	}
	go func() {
		// Launch Gin and
		// handle potential error
		app.Run(ctx)

		wg.Done()
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Gracefully shutting down server...")
	cancel()

	// wait for go routine with api ends shutdown process
	wg.Wait()

	logger.Info("Successfully shutdown")
}

func makeLogger() *zap.SugaredLogger {
	baseLogger, _ := zap.NewProduction()
	//baseLogger, _ := zap.NewDevelopment()
	return baseLogger.Sugar()
}

func checkConfig(cfg AppConfig) error {
	if cfg.Address == "" {
		return errors.New("не указан адрес и порт запуска сервиса")
	}
	if cfg.DatabaseDSN == "" {
		return errors.New("не указан адрес подключения к базе данных")
	}
	if cfg.AccrualSystemAddress == "" {
		return errors.New("не указан адрес системы расчёта начислений")
	}

	if matched, _ := regexp.Match(addressPattern, []byte(cfg.Address)); !matched {
		return errors.New("неверный формат адреса запуска сервиса (host:port)")
	}

	if matched, _ := regexp.Match(addressPattern, []byte(cfg.AccrualSystemAddress)); !matched {
		return errors.New("неверный формат адреса системы расчёта начислений (host:port)")
	}

	return nil
}

func parseFlags(cfg *AppConfig) {
	addressFlag := flag.String("a", "", "адрес и порт запуска сервиса")
	databaseDSNFlag := flag.String("d", "", "адрес подключения к базе данных")
	accrualAddressFlag := flag.String("r", "", "адрес системы расчёта начислений")
	flag.Parse()

	if *addressFlag != "" {
		cfg.Address = *addressFlag
	}
	if *databaseDSNFlag != "" {
		cfg.DatabaseDSN = *databaseDSNFlag
	}
	if *accrualAddressFlag != "" {
		cfg.AccrualSystemAddress = *accrualAddressFlag
	}
}
