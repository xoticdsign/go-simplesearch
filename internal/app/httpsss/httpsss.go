package httpsss

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/xoticdsign/simplesearch/https/simplesearch/ssv1"
	search "github.com/xoticdsign/simplesearch/internal/services/elasticsearch"
	"github.com/xoticdsign/simplesearch/internal/services/simplesearch"
	"github.com/xoticdsign/simplesearch/internal/utils"
)

const op = "app.SimpleSearch."

// App struct represents the SimpleSearch application.
//
// It holds the server and client for handling requests and responses, as well as configuration settings.
type App struct {
	Server ssv1.Server
	Client ssv1.Client

	log    *slog.Logger
	config utils.Config
}

// New initializes and returns a new instance of the SimpleSearch App.
//
// It sets up the Fiber server, initializes the SimpleSearch service, and configures request handlers.
func New(log *slog.Logger, cfg utils.Config) (*App, error) {
	server := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var e *fiber.Error

			if errors.As(err, &e) {
				return c.JSON(&e)
			}
			return c.JSON(fiber.ErrInternalServerError)
		},
		AppName: cfg.ServiceName,
	})

	service, err := simplesearch.New(log, cfg)
	if err != nil {
		return &App{}, err
	}

	handlers := handlers{SimpleSearch: service}

	server.Post("/search", handlers.MakeSearch)

	return &App{
		Server: ssv1.Server{
			ServerImplementation: server,
			Handlers:             handlers,
		},
		Client: ssv1.Client{},

		log:    log,
		config: cfg,
	}, nil
}

// Run starts the SimpleSearch application by listening on the configured host and port.
//
// It listens for incoming HTTP requests and forwards them to the appropriate handler functions.
func (a *App) Run() error {
	const fu = "Run()"

	err := a.Server.ServerImplementation.Listen(utils.BuildAddress(a.config.Host, a.config.Port))
	if err != nil {
		a.log.Warn(
			"listening error",
			slog.String("op", op+fu),
		)

		return err
	}
	return nil
}

// Shutdown shuts down the SimpleSearch application gracefully.
//
// It stops the server and releases any resources held by the application.
func (a *App) Shutdown() error {
	const fu = "Shutdown()"

	err := a.Server.ServerImplementation.Shutdown()
	if err != nil {
		a.log.Warn(
			"shutdown error",
			slog.String("op", op+fu),
		)

		return err
	}
	return nil
}

// Holds all the HTTP request handlers for the SimpleSearch app.
//
// It includes the SimpleSearch service that is responsible for handling search operations.
type handlers struct {
	ssv1.UnimplementedHandlers

	SimpleSearch SimpleSearcher
}

// SimpleSearcher interface defines the contract for searching functionality.
//
// It contains the method `MakeSearch` that takes a search request and returns a list of products or an error.
type SimpleSearcher interface {
	MakeSearch(ctx context.Context, req ssv1.MakeSearchRequest) ([]search.Product, error)
}

// MakeSearch handler processes search requests from clients.
//
// It parses the incoming request, performs validation, delegates the search to the SimpleSearch service,
// and returns the results to the client in JSON format.
func (h *handlers) MakeSearch(c *fiber.Ctx) error {
	var req ssv1.MakeSearchRequest

	err := c.BodyParser(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if req.SearchFor == "" {
		return c.JSON(ssv1.MakeSearchResponse{
			Message: "you must provide the desired search",
		})
	}

	if req.Filters.PriceBottom == 0 && req.Filters.PriceTop == 0 {
		req.Filters.PriceTop = 10000000
	}

	result, err := h.SimpleSearch.MakeSearch(context.Background(), req)
	if err != nil {
		if errors.Is(err, search.ErrNoHits) {
			return c.JSON(ssv1.MakeSearchResponse{
				Message: "none found",
			})
		}
		return fiber.ErrInternalServerError // TODO: BETTER ERROR HANDLING
	}

	return c.JSON(ssv1.MakeSearchResponse{
		Message: "results",
		Result:  result,
	})
}
