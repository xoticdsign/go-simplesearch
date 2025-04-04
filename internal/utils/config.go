package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// App Config Struct.
type Config struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	ServiceName  string        `yaml:"service_name"`

	ElasticSearch ElasticSearch `yaml:"elasticsearch"`
}

// ElasticSearch Struct.
type ElasticSearch struct {
	Addresses []string `yaml:"addresses"`
	Username  string
	Password  string
	Transport ESTransport `yaml:"transport"`
}

// ElasticSearch Transport Struct.
type ESTransport struct {
	TLS         ESTransportTLS `yaml:"tls"`
	TLSTimeout  time.Duration  `yaml:"tls_timeout"`
	IdleTimeout time.Duration  `yaml:"idle_timeout"`
}

// ElasticSearch Transport TLS Struct.
type ESTransportTLS struct {
	Insecure bool `yaml:"tls_insecure"`
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

	esUsername := os.Getenv("ES_USERNAME")
	defer os.Unsetenv("ES_USERNAME")

	esPassword := os.Getenv("ES_PASSWORD")
	defer os.Unsetenv("ES_PASSWORD")

	if esUsername == "" || esPassword == "" {
		return Config{}, fmt.Errorf("following env variables for elasticsearch must be set: ES_USERNAME, ES_PASSWORD")
	}

	cfg.ElasticSearch.Username = esUsername
	cfg.ElasticSearch.Password = esPassword

	return cfg, nil
}
