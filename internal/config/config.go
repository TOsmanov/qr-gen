package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	QRGen      `yaml:"qrGen"`
	HTTPServer `yaml:"httpServer"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"localhost:8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout     time.Duration `yaml:"idleTimeout" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout" env-default:"10s"`
}
type QRGen struct {
	OutputDir   string `yaml:"outputDir" env-default:"output"`
	TempDir     string `yaml:"tempDir" env-default:"temp"`
	SiteDir     string `yaml:"siteDir" env-default:"site"`
	PreviewPath string
	MainPage    string
}

func MustLoad() *Config {
	configPath := getEnv("CONFIG_PATH", "./config/local.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}
	cfg.PreviewPath = fmt.Sprintf("%s/preview.jpg", cfg.SiteDir)
	cfg.MainPage = fmt.Sprintf("%s/index.html", cfg.SiteDir)

	return &cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
