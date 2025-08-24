package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	router "github.com/SLANGERES/go-service/internal/Routers"
)

func TestRoutesExist(t *testing.T) {
	r := router.Router()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/health", nil)
    r.ServeHTTP(w, req)
    if w.Code == http.StatusNotFound {
        t.Errorf("expected route /health to exist")
	}
}
func TestLogRouteExist(t *testing.T) {
	r := router.Router()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/log", nil)
    r.ServeHTTP(w, req)
    if w.Code == http.StatusNotFound {
        t.Errorf("expected route /log to exist")
	}
}