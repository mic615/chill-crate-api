package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mic615/chill-crate-api/internal/models"
)

type NewBucket struct {
	Name    string    `json:"name"     binding:"required"`
	GroupID uuid.UUID `json:"group_id" binding:"required"`
}

func (h *Handler) CreateBucket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bucket NewBucket
		if err := c.ShouldBindJSON(&bucket); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var group models.Group
		if err := h.db.First(&group, "id = ?", bucket.GroupID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		bucketID := uuid.New()
		newBucket := models.Bucket{ID: bucketID, Name: bucket.Name, GroupID: bucket.GroupID}
		if err := h.storageClient.CreateBucket(bucketID.String()); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := h.db.Create(&newBucket).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, newBucket)
	}
}

func (h *Handler) GetBucketsByGroupID() gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("groupId")
		buckets := []models.Bucket{}

		var group models.Group
		if err := h.db.First(&group, "id = ?", groupID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if err := h.db.Where("group_id = ?", groupID).Find(&buckets).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, buckets)
	}
}

func (h *Handler) GetBucketByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("groupId")
		name := c.Param("name")
		bucket := models.Bucket{}

		var group models.Group
		if err := h.db.First(&group, "id = ?", groupID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if err := h.db.Where("group_id = ? AND name = ?", groupID, name).
			First(&bucket).
			Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, bucket)
	}
}

func (h *Handler) DeleteBucket() gin.HandlerFunc {
	return func(c *gin.Context) {
		bucketID := c.Param("bucketId")
		force := strings.ToLower(c.Query("force")) == "true"
		var bucket models.Bucket
		if err := h.db.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}

		// find all objects in the bucket
		objects := []models.Object{}
		if err := h.db.Where("bucket_id = ? AND delete_marker = false", bucketID).
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
			if err := h.storageClient.DeleteObjects(bucket.ID.String(), objects); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		if err := h.storageClient.DeleteBucket(bucket.ID.String()); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// after S3 objects + S3 bucket are gone, clear the DB in a transaction
		err := h.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Unscoped().
				Where("bucket_id = ?", bucket.ID).
				Delete(&models.Object{}).
				Error; err != nil {
				return err
			}
			return tx.Delete(&bucket).Error
		})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, "bucket deleted")
	}
}
