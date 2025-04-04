package logger

import (
	"log/slog"
	"os"
)

func New(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "production":
		// TODO

	case "development":
		// TODO

	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
