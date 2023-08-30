package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App
	HTTP
	Log
	Postgres
}

type App struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type HTTP struct {
	Port string `yaml:"port"`
}

type Log struct {
	LogLevel string `yaml:"log_level"`
}

type Postgres struct {
	PoolMax int    `yaml:"pool_max"`
	URL     string `yaml:"url"`
}

// NewConfig returns app config.
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
