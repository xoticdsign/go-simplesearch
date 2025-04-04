package httpsss

import (
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
func New(log *slog.Logger, cfg utils.Config) *App {
	server := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorHandler: errHandler,
		AppName:      cfg.ServiceName,
	})

	// TODO: IMPLEMENT HANDLERS

	return &App{
		Server: ssv1.Server{
			ServerImplementation: server,
			Service:              simplesearch.New(log),
		},
		Client: ssv1.Client{},

		log:    log,
		config: cfg,
	}
}

func errHandler(c *fiber.Ctx, err error) error {
	return nil
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
