package handlers

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mic615/chill-crate-api/internal/auth"
	"github.com/mic615/chill-crate-api/internal/models"
	"github.com/mic615/chill-crate-api/internal/storage"
)

type Handler struct {
	db            *gorm.DB
	storageClient *storage.Storage
}

func NewHandler(db *gorm.DB, storageClient *storage.Storage) *Handler {
	return &Handler{db: db, storageClient: storageClient}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	v, exists := c.Get("user")
	if !exists {
		return uuid.Nil, false
	}
	user, ok := v.(*models.User)
	if !ok {
		return uuid.Nil, false
	}
	return user.ID, true
}

// returns true if authorized; writes the error response and returns false if not
func (h *Handler) authorize(c *gin.Context, groupID uuid.UUID, role models.Role) bool {
	userID, ok := GetUserID(c)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return false
	}
	if err := auth.RequireRole(h.db, userID, groupID, role); err != nil {
		if errors.Is(err, auth.ErrForbidden) {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return false
	}
	return true
}
