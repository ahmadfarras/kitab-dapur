package platform

import (
	"context"
	"log/slog"
	"time"

	"kitab-dapur/api/internal/openapi"
)

const healthTimeout = 2 * time.Second

type pinger interface {
	Ping(ctx context.Context) error
}

// Health implements the health part of openapi.StrictServerInterface.
type Health struct {
	db  pinger
	log *slog.Logger
}

// NewHealth returns a Health that checks db.
func NewHealth(db pinger, log *slog.Logger) *Health {
	return &Health{db: db, log: log}
}

// GetHealth handles GET /health.
func (h *Health) GetHealth(ctx context.Context, _ openapi.GetHealthRequestObject) (openapi.GetHealthResponseObject, error) {
	ctx, cancel := context.WithTimeout(ctx, healthTimeout)
	defer cancel()
	if err := h.db.Ping(ctx); err != nil {
		h.log.ErrorContext(ctx, "health check: database unreachable", "err", err)
		return openapi.GetHealth503JSONResponse{Message: "Database tidak bisa dihubungi"}, nil
	}
	return openapi.GetHealth200JSONResponse{Status: "ok"}, nil
}
