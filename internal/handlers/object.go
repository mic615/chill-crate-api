package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mic615/chill-crate-api/internal/database"
	"github.com/mic615/chill-crate-api/internal/models"
)

type NewObject struct {
	FileName string `json:"file_name" binding:"required"`
}

func UploadObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var object NewObject
		bucketID := c.Param("bucketId")
		if err := c.ShouldBindJSON(&object); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// todo validate and sanitize file name

		// Todo add role checks when auth is added
		// verify bucket exists
		var bucket models.Bucket
		if err := database.DB.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}
		// calculate version
		objects := []models.Object{}
		database.DB.Where("bucket_id = ? AND file_name = ? ", bucketID, object.FileName).Order("version desc").Find(&objects)
		version := 1
		if len(objects) > 0 {
			version = objects[0].Version + 1
		}
		// todo store file
		newObject := models.Object{FileName: object.FileName, BucketID: bucket.ID, Version: version, Size: 1}
		if err := database.DB.Create(&newObject).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, newObject)
	}

}

func GetObjectsByBucketID() gin.HandlerFunc {
	return func(c *gin.Context) {
		bucketID := c.Param("bucketId")
		objects := []models.Object{}

		var bucket models.Bucket
		if err := database.DB.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}
		if err := database.DB.Where("bucket_id = ?", bucketID).Find(&objects).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(objects) == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no objects found"})
			return
		}
		c.IndentedJSON(http.StatusOK, objects)
	}
}

func GetObjectByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		object := models.Object{}

		if err := database.DB.Find(&object, id).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, object)
	}
}
