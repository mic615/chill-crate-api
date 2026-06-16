package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mic615/chill-crate-api/internal/database"
	"github.com/mic615/chill-crate-api/internal/models"
	"github.com/mic615/chill-crate-api/internal/storage"
)

func UploadObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var object NewObject
		bucketID := c.Param("bucketId")
		// todo validate and sanitize file name
		file, head, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Todo add role checks when auth is added
		// verify bucket exists
		var bucket models.Bucket
		if err := database.DB.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}
		// calculate version
		objects := []models.Object{}
		database.DB.Where("bucket_id = ? AND file_name = ? ", bucketID, head.Filename).Order("version desc").Find(&objects)
		version := 1
		if len(objects) > 0 {
			version = objects[0].Version + 1
		}
		storagePath := uuid.New()
		// todo check file size
		// based on size , do single or multi part upload
		if err := storage.UploadObject(bucket.Name, storagePath.String(), head.Filename, file); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		newObject := models.Object{FileName: head.Filename, BucketID: bucket.ID, Version: version, StoragePath: storagePath, Size: head.Size}
		if err := database.DB.Create(&newObject).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, newObject)
	}

}

func DownloadObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		bucketID := c.Param("bucketId")
		var bucket models.Bucket
		if err := database.DB.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}
		object := models.Object{}
		// get the latest object
		if err := database.DB.Where("bucket_id = ? AND file_name = ? ", bucketID, filename).Order("version desc").First(&object).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "object not found"})
			return
		}
		// todo check file size
		// based on size , do single or multi part downlod
		body, err := storage.DownloadObject(bucket.Name, object.StoragePath.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer body.Close()
		headerMap := map[string]string{"Content-Disposition": fmt.Sprintf("attachment; filename=%q", filename)}

		c.DataFromReader(http.StatusOK, object.Size, "application/octet-stream", body, headerMap)

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
