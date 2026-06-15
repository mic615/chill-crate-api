package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mic615/chill-crate-api/internal/database"
	"github.com/mic615/chill-crate-api/internal/models"
)

type NewGroup struct {
	Name string `json:"name" binding:"required"`
}

func CreateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var group NewGroup
		if err := c.ShouldBindJSON(&group); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Todo add keycloakID when auth is added
		newGroup := models.Group{Name: group.Name, KCGroupID: ""}
		if err := database.DB.Create(&newGroup).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		}
		c.IndentedJSON(http.StatusCreated, newGroup)
	}

}

func GetGroups() gin.HandlerFunc {
	return func(c *gin.Context) {
		groups := []models.Group{}
		if err := database.DB.Find(&groups).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, groups)

	}
}
