package httpsss

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/xoticdsign/simplesearch/https/simplesearch/ssv1"
	"github.com/xoticdsign/simplesearch/internal/services/simplesearch"
	"github.com/xoticdsign/simplesearch/internal/utils"
)

// SimpleSearch App struct.
type App struct {
	Server ssv1.Server
	Client ssv1.Client

	log    *slog.Logger
	config utils.Config
}

// Creates a new SimpleSearch App.
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

// Runs the SimpleSearch App.
func (a *App) Run() error {
	err := a.Server.ServerImplementation.Listen(utils.BuildAddress(a.config.Host, a.config.Port))
	if err != nil {
		return err
	}
	return nil
}

// Shuts down the SimpleSearch App.
func (a *App) Shutdown() error {
	err := a.Server.ServerImplementation.Shutdown()
	if err != nil {
		return err
	}
	return nil
}

// Struct that holds all server Handlers. Also embeeds mocks of handlers for tests.
type handlers struct {
	ssv1.UnimplementedHandlers

	SimpleSearch SimpleSearcher
}

// Interface for SimpleSearch Service.
type SimpleSearcher interface {
	MakeSearch(ctx context.Context, searchFor string) (string, error)
}

// MakeSearch handler.
func (h *handlers) MakeSearch(c *fiber.Ctx) error {
	var req ssv1.MakeSearchRequest

	if req.SearchFor == "" {
		return c.JSON(ssv1.MakeSearchResponse{
			Result: "you must provide the desired search",
		})
	}

	result, err := h.SimpleSearch.MakeSearch(context.Background(), req.SearchFor)
	if err != nil {
		return fiber.ErrInternalServerError // TODO: BETTER ERROR HANDLING
	}

	return c.JSON(ssv1.MakeSearchResponse{
		Result: result,
	})
}
