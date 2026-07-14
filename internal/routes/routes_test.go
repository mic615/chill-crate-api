package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/mic615/chill-crate-api/internal/handlers"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := handlers.NewHandler(nil, nil)
	router := gin.New()
	RegisterRoutes(router, h, func(c *gin.Context) { c.Next() })

	req := httptest.NewRequestWithContext(
		context.Background(), http.MethodGet, "/ping", http.NoBody,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if got := w.Body.String(); got != "OK" {
		t.Fatalf("expected body %q, got %q", "OK", got)
	}
}
