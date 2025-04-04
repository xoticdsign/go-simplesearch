package httpsss

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/xoticdsign/simplesearch/https/simplesearch/ssv1"
	"github.com/xoticdsign/simplesearch/internal/services/simplesearch"
	"github.com/xoticdsign/simplesearch/internal/utils"
)

type App struct {
	Server ssv1.Server
	Client ssv1.Client
}

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
	}
}

func errHandler(c *fiber.Ctx, err error) error {
	return nil
}
