package config

import (
	"server/internal/apperrors"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var path = ".env"

type Config struct {
	AppPort  string `required:"true" split_words:"true"`
	Postgres *PostgresConfig
	Server   *ServerConfig
}

type PostgresConfig struct {
	LogLevel     string `env:"LOGGER_LEVEL"`
	SqlHost      string `env:"SQL_HOST"`
	SqlPort      string `env:"SQL_PORT"`
	SqlType      string `env:"SQL_TYPE"`
	SqlMode      string `env:"SQL_MODE"`
	UserName     string `env:"USER_NAME"`
	Password     string `env:"PASSWORD"`
	DBName       string `env:"DB_NAME"`
	TimeZone     string `env:"TIME_ZONE"`
	TimeoutQuery string `env:"TIMEOUT_QUERY"`
}

type ServerConfig struct {
	AppPort                string `env:"APP_PORT"`
	SecretKey              string `env:"SECRET_KEY"`
	ExpirationJWTInSeconds string `env:"EXPIRATION_JWT_SECONDS"`
	TimeoutContext         string `env:"TIMEOUT_CONTEXT"`
}

func NewConfig(logger *logrus.Logger) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		appErr := apperrors.EnvConfigLoadError.AppendMessage(err)
		return nil, appErr
	}

	confPsql := &PostgresConfig{}
	if err := env.Parse(confPsql); err != nil {
		appErr := apperrors.EnvConfigParseError.AppendMessage(err)
		return nil, appErr
	}

	confServer := &ServerConfig{}
	if err := env.Parse(confServer); err != nil {
		appErr := apperrors.EnvConfigParseError.AppendMessage(err)
		return nil, appErr
	}

	conf := Config{AppPort: confServer.AppPort, Postgres: confPsql, Server: confServer}

	logger.Info("Config has been parsed")
	return &conf, nil
}
