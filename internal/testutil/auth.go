package testutil

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	"github.com/mic615/chill-crate-api/internal/models"
)

func AuthenticateRequest(user *models.User) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user", user)
	return w, c
}
