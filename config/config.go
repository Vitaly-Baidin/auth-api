package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App   `yaml:"app"`
		JWT   `yaml:"jwt"`
		HTTP  `yaml:"http"`
		PG    `yaml:"postgres"`
		Redis `yaml:"redis"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"`
		Version string `env-required:"true" yaml:"version"`
	}

	JWT struct {
		Secret               string `env-required:"true" yaml:"name" env:"JWT_SECRET"`
		AccessExpireInMinute int    `env-required:"true" yaml:"access_expire_in_minute"`
		RefreshExpireInHour  int    `env-required:"true" yaml:"refresh_expire_in_hour"`
	}

	HTTP struct {
		Port                 string `env-required:"true" yaml:"port"`
		ReadTimeoutInSec     int    `yaml:"read_timeout_in_sec" env-default:"15"`
		WriteTimeoutInSec    int    `yaml:"write_timeout_in_sec" env-default:"15"`
		ShutdownTimeoutInSec int    `yaml:"shutdown_timeout_in_sec" env-default:"15"`
	}

	PG struct {
		URL              string `env-required:"true" yaml:"url" env:"PG_URL"`
		PoolMax          int    `env-required:"true" yaml:"pool_max" env-default:"10"`
		ConnAttempts     int    `yaml:"conn_attempts" env-default:"10"`
		ConnTimeoutInSec int    `yaml:"conn_timeout_in_sec" env-default:"1"`
	}

	Redis struct {
		Addresses []string `env-required:"true" yaml:"addresses"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
