package config

import (
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database struct {
		URI         string `yaml:"postgressql"`
		MaxAttempts int    `yaml:"max-attempts"`
	} `yaml:"database"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	TokenConfig struct {
		SecretKey           string        `yaml:"secret-key"`
		AccessTimeLiveToken time.Duration `yaml:"access-time-live-token""`
	} `yaml:"token"`
}

func NewConfig() *Config {
	var cfg Config
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	path := basePath + "/config.yml"
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatal("config read error:", err)
	}

	return &cfg
}
