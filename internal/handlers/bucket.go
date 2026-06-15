package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mic615/chill-crate-api/internal/database"
	"github.com/mic615/chill-crate-api/internal/models"
)

type NewBucket struct {
	Name    string    `json:"name" binding:"required"`
	GroupID uuid.UUID `json:"group_id"  binding:"required"`
}

func CreateBucket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bucket NewBucket
		if err := c.ShouldBindJSON(&bucket); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var group models.Group
		if err := database.DB.First(&group, "id = ?", bucket.GroupID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		newBucket := models.Bucket{Name: bucket.Name, GroupID: bucket.GroupID}
		if err := database.DB.Create(&newBucket).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, newBucket)
	}
}

func GetBucketsByGroupID() gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("groupId")
		buckets := []models.Bucket{}

		var group models.Group
		if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if err := database.DB.Where("group_id = ?", groupID).Find(&buckets).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(buckets) == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no buckets found"})
			return
		}
		c.IndentedJSON(http.StatusOK, buckets)

	}
}
