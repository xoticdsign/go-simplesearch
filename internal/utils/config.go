package utils

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config struct represents the configuration of the application.
//
// It holds all necessary configuration parameters such as the host, port, timeouts,
// service name, and ElasticSearch-related settings.
type Config struct {
	Address      string
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	ServiceName  string        `yaml:"service_name"`

	ElasticSearch ElasticSearch `yaml:"elasticsearch"`
}

// ElasticSearch struct represents the ElasticSearch connection settings.
//
// It contains the addresses, credentials, and transport-related settings for connecting to ElasticSearch.
type ElasticSearch struct {
	Address   string
	Username  string
	Password  string
	Transport ESTransport `yaml:"transport"`
}

// ESTransport struct represents the transport layer settings for ElasticSearch.
//
// It includes TLS settings and timeouts for the transport layer.
type ESTransport struct {
	TLS         ESTransportTLS `yaml:"tls"`
	TLSTimeout  time.Duration  `yaml:"tls_timeout"`
	IdleTimeout time.Duration  `yaml:"idle_timeout"`
}

// ESTransportTLS struct holds the TLS settings for ElasticSearch.
//
// It specifies whether the connection should ignore insecure TLS settings.
type ESTransportTLS struct {
	Insecure bool `yaml:"tls_insecure"`
}

// MustLoadConfig() loads the application configuration from a .yaml file.
//
// It takes the environment as an argument to determine the config file to use. If the environment is "production"
// or "development", it will load respective config files (not implemented yet). If the environment is not recognized,
// it loads a local config file. The function also retrieves ElasticSearch credentials from environment variables,
// and sets them into the configuration. If any error occurs during the loading process, it returns an error.
func MustLoadConfig(envs map[string]string) (Config, error) {
	var cfg Config

	switch envs["env"] {
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

	cfg.Address = envs["address"]
	cfg.ElasticSearch.Address = envs["es_address"]
	cfg.ElasticSearch.Username = envs["es_username"]
	cfg.ElasticSearch.Password = envs["es_password"]

	return cfg, nil
}
