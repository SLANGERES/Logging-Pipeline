package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	router "github.com/SLANGERES/go-service/internal/Routers"
)

func TestMainStartup(t *testing.T) {
	// Create router directly (don’t bind to port)
	r := router.Router()

	// Use httptest server (no need to run r.Run())
	ts := httptest.NewServer(r)

	defer ts.Close()

	// Call /health endpoint
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("could not GET /health: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}
