package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/xoticdsign/simplesearch/internal/app/httpsss"
	"github.com/xoticdsign/simplesearch/internal/lib/logger"
	"github.com/xoticdsign/simplesearch/internal/utils"
)

// App Struct.
type App struct {
	SimpleSearch *httpsss.App

	log    *slog.Logger
	config utils.Config
}

// Creates a new App.
//
// Returns App Struct, if everything's ok. Returns Error, if something went wrong while creating one of the components.
func New(env string) (*App, error) {
	cfg, err := utils.MustLoadConfig(env)
	if err != nil {
		return &App{}, err
	}

	log := logger.New(env)

	ss, err := httpsss.New(log, cfg)
	if err != nil {
		return &App{}, err
	}

	return &App{
		SimpleSearch: ss,

		log:    log,
		config: cfg,
	}, nil
}

// Runs the App.
//
// Each part of the App starts in a separate goroutine and then Run function waits for SIGTERM/SIGINT or an Error to occur to proceed. After Signal or an Error has been received, Run function call for shutdown and shuts down gracefully.
func (a *App) Run() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error, 1)

	go func() {
		err := a.SimpleSearch.Run()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		a.log.Error(
			"will shutdown, because an error occurred while running simplesearch server",
			slog.String("error", err.Error()),
		)

	case <-sigChan:
		a.log.Info(
			"signaled to shutdown",
		)
	}

	err := a.shutdown()
	if err != nil {
		return err
	}
	return nil
}

// Shuts down the App.
func (a *App) shutdown() error {
	err := a.SimpleSearch.Shutdown()
	if err != nil {
		a.log.Error(
			"error occurred while trying to shutdown gracefully",
			slog.String("error", err.Error()),
		)

		return err
	}
	return nil
}
