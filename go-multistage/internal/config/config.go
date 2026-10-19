package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	App    AppConfig    `yaml:"app"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type AppConfig struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
}

func Load() (*Config, error) {
	configFile := getEnv("CONFIG_FILE", "config.yaml")

	cfg := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		App: AppConfig{
			Name:        "go-app",
			Environment: "development",
		},
	}

	if err := loadFile(configFile, cfg); err != nil {
		return nil, err
	}

	if err := loadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadFile(filename string, cfg *Config) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Config file is optional.
			return nil
		}

		return fmt.Errorf("read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}

func loadEnv(cfg *Config) error {
	if value := os.Getenv("SERVER_HOST"); value != "" {
		cfg.Server.Host = value
	}

	if value := os.Getenv("SERVER_PORT"); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid SERVER_PORT %q: %w", value, err)
		}

		cfg.Server.Port = port
	}

	if value := os.Getenv("APP_NAME"); value != "" {
		cfg.App.Name = value
	}

	if value := os.Getenv("APP_ENVIRONMENT"); value != "" {
		cfg.App.Environment = value
	}

	return nil
}

func getEnv(name, fallback string) string {
	value := os.Getenv(name)

	if value == "" {
		return fallback
	}

	return value
}
