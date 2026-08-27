package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stubPinger struct {
	err error
}

func (s stubPinger) Ping(context.Context) error {
	return s.err
}

func TestHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pingErr    error
		wantStatus int
		want       healthResponse
	}{
		{
			name:       "ok",
			pingErr:    nil,
			wantStatus: http.StatusOK,
			want:       healthResponse{Status: "ok", DB: "ok"},
		},
		{
			name:       "fail",
			pingErr:    errors.New("connection refused"),
			wantStatus: http.StatusServiceUnavailable,
			want:       healthResponse{Status: "degraded", DB: "error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := NewMux(stubPinger{err: tt.pingErr})
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var got healthResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode body: %v\nbody: %q", err, rec.Body.String())
			}
			if got != tt.want {
				t.Fatalf("body = %+v, want %+v", got, tt.want)
			}
			if tt.pingErr != nil && strings.Contains(rec.Body.String(), tt.pingErr.Error()) {
				t.Fatalf("response leaked ping error: %q", rec.Body.String())
			}
		})
	}
}
