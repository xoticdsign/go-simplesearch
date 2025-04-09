package simplesearch

import (
	"context"
	"log/slog"

	"github.com/xoticdsign/go-simplesearch/https/simplesearch/ssv1"
	search "github.com/xoticdsign/go-simplesearch/internal/services/elasticsearch"
	"github.com/xoticdsign/go-simplesearch/internal/utils"
)

const op = "service.SimpleSearch."

// Service struct represents the SimpleSearch service.
//
// It contains the necessary dependencies such as the Searcher (interface for search engines),
// logger for logging, and configuration settings.
type Service struct {
	Search Searcher

	log    *slog.Logger
	config utils.Config
}

// Searcher interface defines the contract for search engines used by the SimpleSearch service.
//
// It has a method 'MakeSearch' that takes a request object and returns search results (as Product structs) and an error.
type Searcher interface {
	MakeSearch(req ssv1.MakeSearchRequest) ([]search.Product, error)
}

// New initializes and returns a new instance of the SimpleSearch service.
//
// It creates a new search engine client (such as Elasticsearch) and passes the logger and configuration settings.
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

// MakeSearch is a method on the SimpleSearch service that performs a search using the provided request.
//
// It delegates the search operation to the underlying Searcher interface (e.g., Elasticsearch client).
// If the search operation is successful, it returns the results, otherwise it returns an error.
func (s *Service) MakeSearch(ctx context.Context, req ssv1.MakeSearchRequest) ([]search.Product, error) {
	const fu = "MakeSearch()"

	result, err := s.Search.MakeSearch(req)
	if err != nil {
		return []search.Product{}, err
	}
	return result, nil
}
