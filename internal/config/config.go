package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL,required,unset"`
	JWTSecret   string `env:"JWT_SECRET,required,unset"`
	Port        uint16 `env:"PORT" envDefault:"8080"`
}

func Load() (config *Config, err error) {
	config = new(Config)
	if err = env.Parse(config); err != nil {
		return
	}
	return
}
