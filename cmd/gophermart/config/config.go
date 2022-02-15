package config

import (
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"regexp"
)

var addressPattern = `^.+?:\d{1,5}$`

type Config struct {
	Address              string `env:"RUN_ADDRESS" envDefault:"localhost:8000"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseDSN          string `env:"DATABASE_URI"`
	MigrationsDir        string `env:"DATABASE_MIGRATIONS" envDefault:"migrations"`
}

func Parse() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)

	if err != nil {
		return nil, err
	}
	cfg.parseFlags()
	err = cfg.checkConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) checkConfig() error {
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

func (cfg *Config) parseFlags() {
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
