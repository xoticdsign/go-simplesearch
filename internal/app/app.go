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

const op = "app."

// App Struct represents the entire application.
//
// It contains the SimpleSearch service and manages the application's lifecycle (initialization, running, and shutdown).
type App struct {
	SimpleSearch *httpsss.App

	log    *slog.Logger
	config utils.Config
}

func getEnv() (map[string]string, error) {
	env := os.Getenv("ENV")
	address := os.Getenv("ADDRESS")
	esAddress := os.Getenv("ES_ADDRESS")
	esUsername := os.Getenv("ES_USERNAME")
	esPassword := os.Getenv("ES_PASSWORD")

	envs := make(map[string]string)

	envs["env"] = env
	envs["address"] = address
	envs["es_address"] = esAddress
	envs["es_username"] = esUsername
	envs["es_password"] = esPassword

	return envs, nil
}

// New() creates a new instance of the App.
//
// It loads the configuration based on the provided environment, initializes the logger,
// and creates the SimpleSearch application. Returns the App struct or an error if any step fails.
func New() (*App, error) {
	envs, err := getEnv()
	if err != nil {
		return &App{}, err
	}

	cfg, err := utils.MustLoadConfig(envs)
	if err != nil {
		return &App{}, err
	}

	log := logger.New(envs)

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

// Run() starts the application and handles the main application flow.
//
// It starts the SimpleSearch app in a separate goroutine and listens for termination signals (SIGTERM, SIGINT).
// If an error occurs or a shutdown signal is received, it proceeds to shut down the application gracefully.
func (a *App) Run() error {
	const fu = "Run()"

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
			"will shutdown, because an error occurred while running",
			slog.String("op", op+fu),
			slog.String("error", err.Error()),
		)

	case <-sigChan:
		a.log.Info(
			"signaled to shutdown",
			slog.String("op", op+fu),
		)
	}

	err := a.shutdown()
	if err != nil {
		return err
	}
	return nil
}

// Shuts down the application gracefully.
//
// It stops the SimpleSearch service and handles any errors that may occur during the shutdown process.
func (a *App) shutdown() error {
	const fu = "shutdown()"

	err := a.SimpleSearch.Shutdown()
	if err != nil {
		a.log.Error(
			"error occurred while trying to shutdown gracefully",
			slog.String("op", op+fu),
			slog.String("error", err.Error()),
		)

		return err
	}
	return nil
}
