package elasticsearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/xoticdsign/simplesearch/https/simplesearch/ssv1"
	"github.com/xoticdsign/simplesearch/internal/utils"
)

const op = "service.ElasticSearch."

var (
	ErrNoHits              = fmt.Errorf("no hits")
	ErrEncodingJSON        = fmt.Errorf("json encoding error")
	ErrDecodingJSON        = fmt.Errorf("json decoding error")
	ErrMarshalingJSON      = fmt.Errorf("json marshaling error")
	ErrUnmarshalingJSON    = fmt.Errorf("json umarshaling error")
	ErrInterfaceConversion = fmt.Errorf("interface conversion error")
)

// Constant representing the ElasticSearch Products index name.
const iProducts = "products"

// Product struct represents a product with its properties that will be unmarshaled from Elasticsearch.
type Product struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}

var (
	pName        = "name^3"
	pDescription = "description^2"
	pPrice       = "price"
	pCategory    = "category^2"
	pStock       = "stock"
	pCreatedAt   = "created_at"
)

// Service struct represents the Elasticsearch service with the necessary client and configurations.
type Service struct {
	ESClient *elasticsearch.Client

	log    *slog.Logger
	config utils.Config
}

// New creates a new instance of the ElasticSearch Service with the given logger and configuration.
//
// It initializes the ElasticSearch client with TLS and authentication settings.
func New(log *slog.Logger, cfg utils.Config) (*Service, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ElasticSearch.Address},
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

// Takes an interface{} (typically the decoded Elasticsearch response),
// extracts the 'hits' from the response, and unmarshals them into a slice of Product structs.
// It returns a slice of Product structs and an error, or an empty slice and a specific error
// if the conversion or unmarshaling process fails.
//
// The function expects the input to be in a format that contains a 'hits' field in the response,
// which is a common structure in Elasticsearch search results. Each hit is a product document
// that is converted to a Product struct. If the structure doesn't match, or there are issues
// with marshaling or unmarshaling, an appropriate error is returned.
func (s *Service) productHitsExtractor(v any) ([]Product, error) {
	const fu = "producHitsEctractor()"

	hits, ok := v.(map[string]interface{})["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		s.log.Error(
			"hits conversion error",
			slog.String("op", op+fu),
			slog.String("error", ErrInterfaceConversion.Error()),
		)

		return []Product{}, ErrInterfaceConversion
	}
	var products []Product

	for _, hit := range hits {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if !ok {
			s.log.Error(
				"hit conversion error",
				slog.String("op", op+fu),
				slog.String("error", ErrInterfaceConversion.Error()),
			)

			return []Product{}, ErrInterfaceConversion
		}
		productJSON, err := json.Marshal(source)
		if err != nil {
			s.log.Error(
				"can't marshal a hit",
				slog.String("op", op+fu),
				slog.String("error", err.Error()),
			)

			return []Product{}, ErrMarshalingJSON
		}

		var product Product

		err = json.Unmarshal(productJSON, &product)
		if err != nil {
			s.log.Error(
				"can't umarshal a hit",
				slog.String("op", op+fu),
				slog.String("error", err.Error()),
			)

			return []Product{}, ErrUnmarshalingJSON
		}

		products = append(products, product)
	}

	return products, nil
}

// MakeSearch performs a search query against Elasticsearch with the given request parameters.
//
// It uses the query, filters, and other parameters specified in the 'req' argument.
func (s *Service) MakeSearch(req ssv1.MakeSearchRequest) ([]Product, error) {
	const fu = "Search()"

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  req.SearchFor,
							"fields": []string{pName},
						},
					},
				},
				"must_not": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							pStock: 0,
						},
					},
				},
				"filter": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							pPrice: map[string]interface{}{
								"gte": req.Filters.PriceBottom,
								"lte": req.Filters.PriceTop,
							},
						},
					},
				},
			},
		},
	}

	buf, err := utils.JSONEncode(query)
	if err != nil {
		s.log.Error(
			"can't encode",
			slog.String("op", op+fu),
			slog.String("error", err.Error()),
		)

		return []Product{}, ErrEncodingJSON
	}

	resp, err := s.ESClient.Search(
		s.ESClient.Search.WithContext(context.Background()),
		s.ESClient.Search.WithIndex(iProducts),
		s.ESClient.Search.WithBody(&buf),
		s.ESClient.Search.WithPretty(),
	)
	if err != nil {
		s.log.Error(
			"can't make a request to elasticsearch",
			slog.String("op", op+fu),
			slog.String("error", err.Error()),
		)

		return []Product{}, err
	}
	defer resp.Body.Close()

	r, err := utils.JSONDecode(resp.Body)
	if err != nil {
		s.log.Error(
			"can't decode",
			slog.String("op", op+fu),
			slog.String("error", err.Error()),
		)

		return []Product{}, ErrDecodingJSON
	}

	products, err := s.productHitsExtractor(r)
	if err != nil {
		return []Product{}, err
	}

	if len(products) == 0 {
		return []Product{}, ErrNoHits
	}

	return products, nil
}
