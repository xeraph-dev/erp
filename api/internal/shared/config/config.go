package config

import (
	"log/slog"
	"os"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	DatabaseURL     string `env:"DATABASE_URL,required,unset"`
	JWTSecret       string `env:"JWT_SECRET,required,unset"`
	Port            uint16 `env:"PORT" envDefault:"8080"`
	LogLevelString  string `env:"LOG_LEVEL" envDefault:"info"`
	LogFormatString string `env:"LOG_FORMAT" envDefault:"json"`
}

func Load() (*Config, error) {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (config Config) LogLevel() slog.Level {
	switch config.LogLevelString {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (config Config) LogHandler() slog.Handler {
	opts := &slog.HandlerOptions{Level: config.LogLevel()}

	switch config.LogFormatString {
	case "text":
		return slog.NewTextHandler(os.Stdout, opts)
	default:
		return slog.NewJSONHandler(os.Stdout, opts)
	}
}
