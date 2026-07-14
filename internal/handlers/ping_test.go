package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/mic615/chill-crate-api/internal/handlers"
)

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequestWithContext(
		context.Background(), http.MethodGet, "/ping", http.NoBody,
	)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handlers.Ping()(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if got := w.Body.String(); got != "OK" {
		t.Fatalf("expected body %q, got %q", "OK", got)
	}
}
