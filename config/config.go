package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	Server          ServerConfig
	Database        DatabaseConfig
	Auth            AuthConfig
	ExchangeService ExchangeService
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Test     bool
}

type AuthConfig struct {
	SecretKey string
}

type ExchangeService struct {
	Host string
	Port string
}

func LoadConfig() *Config {
	configFile := flag.String("c", "", "path to config file")
	flag.Parse()

	if *configFile == "" {
		zap.L().Fatal("path to config file is required")
	}

	if _, err := os.Stat(*configFile); err == nil {
		if err = godotenv.Load(*configFile); err != nil {
			zap.L().Fatal(fmt.Sprintf("error loading config file %s: %v", *configFile, err))
		}
	} else {
		zap.L().Error(fmt.Sprintf("config file not found: %s", *configFile))
	}

	cfg := &Config{}

	cfg.Server.Port = os.Getenv("SERVER_PORT")

	cfg.Database.Host = os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	if dbPort != "" {
		if port, err := strconv.Atoi(dbPort); err == nil {
			cfg.Database.Port = port
		} else {
			zap.L().Fatal(fmt.Sprintf("invalid DB_PORT value: %s", dbPort))
		}
	}
	cfg.Database.User = os.Getenv("DATABASE_USER")
	cfg.Database.Password = os.Getenv("DATABASE_PASSWORD")
	cfg.Database.Name = os.Getenv("DATABASE_NAME")

	cfg.Auth.SecretKey = os.Getenv("JWT_KEY")

	cfg.ExchangeService.Host = os.Getenv("EXCHANGE_SERVICE_HOST")
	cfg.ExchangeService.Port = os.Getenv("EXCHANGE_SERVICE_PORT")

	return cfg
}
