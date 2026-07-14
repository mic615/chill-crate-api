package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mic615/chill-crate-api/internal/models"
)

type NewGroup struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) CreateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var group NewGroup
		user, exists := c.Get("user")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		userID := user.(*models.User).ID
		// TODO handle missing ID
		if err := c.ShouldBindJSON(&group); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newGroup := models.Group{Name: group.Name}
		err := h.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&newGroup).Error; err != nil {
				return err
			}
			membership := models.Membership{
				UserID:  userID,
				GroupID: newGroup.ID,
				Role:    models.RoleAdmin,
			}
			return tx.Create(&membership).Error
		})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, newGroup)
	}
}

func (h *Handler) GetMyGroups() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		userID := user.(*models.User).ID
		// TODO handle missing ID
		groups := []models.Group{}
		if err := h.db.Joins("JOIN memberships ON memberships.group_id = groups.id").
			Where("memberships.user_id = ?", userID).
			Find(&groups).
			Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, groups)
	}
}

func (h *Handler) GetGroups() gin.HandlerFunc {
	return func(c *gin.Context) {
		groups := []models.Group{}
		if err := h.db.Find(&groups).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, groups)
	}
}
