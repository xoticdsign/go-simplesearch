package simplesearch

import (
	"log/slog"

	search "github.com/xoticdsign/simplesearch/internal/services/elasticsearch"
	"github.com/xoticdsign/simplesearch/internal/utils"
)

// SimpleSearch Service Struct.
type Service struct {
	Search Searcher

	log    *slog.Logger
	config utils.Config
}

// Searcher interface.
//
// Responsible for providing methods to use Search Engines.
type Searcher interface{}

// Creates a new SimpleSearch Service.
func New(log *slog.Logger, cfg utils.Config) (*Service, error) {
	search, err := search.New(log, cfg)
	if err != nil {
		return &Service{}, err
	}

	return &Service{
		Search: search,

		log:    log,
		config: cfg,
	}, nil
}
