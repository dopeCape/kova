package config

import (
	util "github.com/dopeCape/kova/internal/utils"
)

type Env string

const (
	DEVELOPMENT Env = "development"
	PRODUCTION  Env = "production"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  Env
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type AuthConfig struct {
	JWTSecret string
}

func Load() *Config {
	env := Env(util.GetEnv("ENVIRONMENT", string(DEVELOPMENT)))
	return &Config{
		Server: ServerConfig{
			Port: util.GetEnv("PORT", "8000"),
			Host: util.GetEnv("HOST", "localhost"),
			Env:  env,
		},
		Database: DatabaseConfig{
			Host:     util.GetEnv("DB_HOST", "localhost"),
			Port:     util.GetEnvInt("DB_PORT", 5432),
			User:     util.GetEnv("DB_USER", "admin"),
			Password: util.GetEnv("DB_PASSWORD", "password123"),
			Database: util.GetEnv("DB_NAME", "mydb"),
			SSLMode:  util.GetEnv("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret: util.GetEnv("JWT_SECRET", "849cff22c983fb7a0ee113339c6486893c83f1e5d485ef2a797b43f802b21709"),
		},
	}
}
