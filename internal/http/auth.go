package http

import (
	"context"
	"encoding/json"
	"net/http"
)

type Auther interface {
	Login(ctx context.Context, username, password string) (string, error)
	Signup(ctx context.Context, username, password string) (string, error)
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Error string `json:"error,omitempty"`
}

type AuthHandler struct {
	DB Auther
}

func (a AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthPingTimeout)
	defer cancel()

	// check pool
	if a.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, AuthResponse{
			Error: "Database connection is not available",
		})
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{
			Error: "Invalid request",
		})
		return
	}

	id, err := a.DB.Login(ctx, req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, AuthResponse{
			Error: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, AuthResponse{
		ID:    id,
		Email: req.Email,
	})
}

func (a AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), healthPingTimeout)
	defer cancel()

	// check pool
	if a.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, AuthResponse{
			Error: "Database connection is not available",
		})
		return
	}
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{
			Error: "Invalid request",
		})
		return
	}

	id, err := a.DB.Signup(ctx, req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, AuthResponse{
			Email: req.Email,
			Error: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, AuthResponse{
		ID:    id,
		Email: req.Email,
	})
}
