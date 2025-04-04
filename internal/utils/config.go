package utils

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// App Config struct.
type Config struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	ServiceName  string        `yaml:"service_name"`
}

// Loads Config from .yaml file.
//
// Environment must be specified with env variable. Return Config struct, everything is ok. Returns Error, if something goes wrong.
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
