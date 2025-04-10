package logger

import (
	"log/slog"
	"os"
)

// Creates a new Logger.
//
// Environment must be specified with env variable.
func New(envs map[string]string) *slog.Logger {
	var log *slog.Logger

	switch envs["env"] {
	case "production":
		// TODO

	case "development":
		// TODO

	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
