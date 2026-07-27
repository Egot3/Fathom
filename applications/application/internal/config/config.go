package config

import (
	"os"
	"slices"
	"strings"
)

/* type Config struct {
	Database struct {
		Sqlite struct {
			Used bool   `yaml:"used"`
			Path string `yaml:"path"`
		} `yaml:"sqlite"`
		Postgres struct {
			Used     bool   `yaml:"used"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			DbName   string `yaml:"dbName"`
		} `yaml:"postgres"`
	} `yaml:"database"`
	Server struct {
		Port    string `yaml:"port"`
		Logging struct {
			Loggers []string `yaml:"loggers"`
			Level   string   `yaml:"level"`
		} `yaml:"logging"`
	} `yaml:"server"`
	QuizPath string `yaml:"quizPath"`
} */

type Config struct {
	DatabaseDriver string
	SqlitePath     string
	PostgresURL    string

	ServerPort string
	LogLevel   string
	LogSinks   []string

	InitAdminUsername string
	InitAdminPassword string

	QuizPath string
}

func getEnv(key, fallback string, allowed []string) string {
	if v, ok := os.LookupEnv(key); ok && (allowed == nil || slices.Contains(allowed, v)) {
		return v

	}
	return fallback
}

func Load() Config {
	pathToQuizzes = "/data/quizzes"
	return Config{
		DatabaseDriver: getEnv("DATABASE_DRIVER", "sqlite", []string{"sqlite", "postgres"}),
		SqlitePath:     "/data/fathom.db",
		PostgresURL:    os.Getenv("POSTGRES_URL"),

		ServerPort: getEnv("SERVER_PORT", "8080", nil),
		LogLevel:   getEnv("LOG_LEVEL", "info", nil),
		LogSinks:   strings.Split(os.Getenv("LOG_SINKS"), ","),

		InitAdminUsername: os.Getenv("INIT_ADMIN_USERNAME"),
		InitAdminPassword: os.Getenv("INIT_ADMIN_PASSWORD"),

		QuizPath: pathToQuizzes,
	}
}
