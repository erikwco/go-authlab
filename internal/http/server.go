package http

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
)

// NewServer returns an HTTP server with timeouts. pinger is used by GET /health.
func NewServer(addr string, pinger Pinger) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           NewMux(pinger),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}

// NewMux registers HTTP routes.
func NewMux(pinger Pinger) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /health", HealthHandler{DB: pinger})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
