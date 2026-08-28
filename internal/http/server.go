// Package http provides HTTP server and handlers.
package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/erikwco/go-authlab/internal/db"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
)

// NewServer returns an HTTP server with timeouts. pinger is used by GET /health.
func NewServer(addr string, store *db.Pool) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           NewMux(store),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}

// NewMux registers HTTP routes.
func NewMux(store *db.Pool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", AuthHandler{DB: store}.Login)
	mux.HandleFunc("POST /signup", AuthHandler{DB: store}.Signup)
	mux.Handle("GET /health", HealthHandler{DB: store})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
