package simplesearch

import "log/slog"

// SimpleSearch Service Struct.
type Service struct {
	log *slog.Logger
}

// Creates a new SimpleSearch Service.
func New(log *slog.Logger) *Service {
	return &Service{
		log: log,
	}
}
