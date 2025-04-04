package elasticsearch

import (
	"crypto/tls"
	"log/slog"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/xoticdsign/simplesearch/internal/utils"
)

// ElasticSearch Service Struct.
type Service struct {
	ESClient *elasticsearch.Client

	log    *slog.Logger
	config utils.Config
}

// Creates a new ElasticSearch Service.
func New(log *slog.Logger, cfg utils.Config) (*Service, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.ElasticSearch.Addresses,
		Username:  cfg.ElasticSearch.Username,
		Password:  cfg.ElasticSearch.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.ElasticSearch.Transport.TLS.Insecure,
			},
			TLSHandshakeTimeout: cfg.ElasticSearch.Transport.TLSTimeout,
			IdleConnTimeout:     cfg.ElasticSearch.Transport.IdleTimeout,
		},
	})
	if err != nil {
		return &Service{}, err
	}
	return &Service{
		ESClient: es,

		log:    log,
		config: cfg,
	}, nil
}
