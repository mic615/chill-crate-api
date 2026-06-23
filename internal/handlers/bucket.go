package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mic615/chill-crate-api/internal/database"
	"github.com/mic615/chill-crate-api/internal/models"
	"github.com/mic615/chill-crate-api/internal/storage"
)

type NewBucket struct {
	Name    string    `json:"name"     binding:"required"`
	GroupID uuid.UUID `json:"group_id" binding:"required"`
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
		bucketID := uuid.New()
		newBucket := models.Bucket{ID: bucketID, Name: bucket.Name, GroupID: bucket.GroupID}
		if err := storage.CreateBucket(bucketID.String()); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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
		c.IndentedJSON(http.StatusOK, buckets)
	}
}

func GetBucketByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("groupId")
		name := c.Param("name")
		bucket := models.Bucket{}

		var group models.Group
		if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if err := database.DB.Where("group_id = ? AND name = ?", groupID, name).
			First(&bucket).
			Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, bucket)
	}
}

func DeleteBucket() gin.HandlerFunc {
	return func(c *gin.Context) {
		bucketID := c.Param("bucketId")
		force := false
		if strings.ToLower(c.Query("force")) == "true" {
			force = true
		}
		var bucket models.Bucket
		if err := database.DB.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}

		// find all objects in the bucket
		objects := []models.Object{}
		if err := database.DB.Where("bucket_id = ? AND delete_marker = false", bucketID).
			Find(&objects).
			Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(objects) > 0 {
			if !force {
				c.IndentedJSON(http.StatusConflict, gin.H{"error": "bucket not empty"})
				return
			}
			err := storage.DeleteObjects(bucket.ID.String(), objects)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
				return

			}
		}
		if err := storage.DeleteBucket(bucket.ID.String()); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if err := database.DB.Delete(&bucket).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, "bucket deleted")
	}
}
