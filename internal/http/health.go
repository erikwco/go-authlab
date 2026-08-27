package http

import (
	"context"
	"net/http"
	"time"
)

const healthPingTimeout = 2 * time.Second

// Pinger is the DB dependency for health checks. Tests stub this.
type Pinger interface {
	Ping(ctx context.Context) error
}

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}

// HealthHandler serves GET /health. Failures never include error text or traces.
type HealthHandler struct {
	DB Pinger
}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthPingTimeout)
	defer cancel()

	if h.DB == nil || h.DB.Ping(ctx) != nil {
		writeJSON(w, http.StatusServiceUnavailable, healthResponse{
			Status: "degraded",
			DB:     "error",
		})
		return
	}

	writeJSON(w, http.StatusOK, healthResponse{
		Status: "ok",
		DB:     "ok",
	})
}
