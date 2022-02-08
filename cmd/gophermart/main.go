package main

import (
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type EnvConfig struct {
	Address              string `env:"RUN_ADDRESS"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseDSN          string `env:"DATABASE_URI"`
}

func main() {
	cfg := EnvConfig{
		Address:              "localhost:8000",
		AccrualSystemAddress: "",
		DatabaseDSN:          "",
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

}

func checkConfig(cfg EnvConfig) error {
	if cfg.Address == "" {
		return errors.New("не указан адрес и порт запуска сервиса")
	}
	if cfg.DatabaseDSN == "" {
		return errors.New("не указан адрес подключения к базе данных")
	}
	if cfg.AccrualSystemAddress == "" {
		return errors.New("не указан адрес системы расчёта начислений")
	}
	return nil
}

func parseFlags(cfg *EnvConfig) {
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
