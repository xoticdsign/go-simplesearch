package utils

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	ServiceName  string        `yaml:"service_name"`
}

func MustLoadConfig(env string) (Config, error) {
	var cfg Config

	switch env {
	case "production":
		// TODO

	case "development":
		// TODO

	default:
		err := cleanenv.ReadConfig("./config/local.yaml", &cfg)
		if err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}
