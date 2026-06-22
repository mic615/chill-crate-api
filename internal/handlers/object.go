package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mic615/chill-crate-api/internal/database"
	"github.com/mic615/chill-crate-api/internal/models"
	"github.com/mic615/chill-crate-api/internal/storage"
)

func UploadObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		bucketID := c.Param("bucketId")
		fileName := c.Param("filename")
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
		database.DB.Where("bucket_id = ? AND file_name = ? ", bucketID, fileName).
			Order("version desc").
			Find(&objects)
		version := 1
		if len(objects) > 0 {
			version = objects[0].Version + 1
		}
		storagePath := uuid.New()
		// todo check file size
		// based on size , do single or multi part upload
		// load small files into buffer
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "reading body: " + err.Error()})
			return
		}

		if err := storage.UploadObject(
			bucket.Name,
			storagePath.String(),
			fileName,
			bytes.NewReader(data),
			int64(len(data)),
		); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// TODO stream large files
		newObject := models.Object{
			FileName:    fileName,
			BucketID:    bucket.ID,
			Version:     version,
			StoragePath: &storagePath,
			Size:        c.Request.ContentLength,
		}
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
		if err := database.DB.Where("bucket_id = ? AND file_name = ?", bucketID, filename).
			Order("version desc").
			First(&object).
			Error; err != nil || object.DeleteMarker {
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
		headerMap := map[string]string{
			"Content-Disposition": fmt.Sprintf("attachment; filename=%q", filename),
		}

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
		queryString := `
			SELECT * FROM (
			SELECT DISTINCT ON (file_name) * 
			FROM objects 
			WHERE bucket_id = ? AND deleted_at IS NULL 
			ORDER BY file_name, version DESC
			) latest 
			WHERE delete_marker = false;
		`
		if err := database.DB.Raw(queryString, bucketID).Scan(&objects).Error; err != nil {
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

func DeleteObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		bucketID := c.Param("bucketId")
		fileName := c.Param("filename")
		// todo validate and sanitize file name

		// Todo add role checks when auth is added
		// verify bucket exists
		var bucket models.Bucket
		if err := database.DB.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}
		// find latest object version and verify it is not deleted already
		object := models.Object{}
		if err := database.DB.Where("bucket_id = ? AND file_name = ? ", bucketID, fileName).
			Order("version desc").
			First(&object).Error; err != nil || object.DeleteMarker {
			c.JSON(http.StatusNotFound, gin.H{"error": "object not found"})
			return
		}
		// calculate version
		version := object.Version + 1
		newObject := models.Object{
			FileName:     fileName,
			BucketID:     bucket.ID,
			Version:      version,
			StoragePath:  nil,
			DeleteMarker: true,
			Size:         0,
		}
		if err := database.DB.Create(&newObject).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, newObject)
	}
}

func RestoreObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		bucketID := c.Param("bucketId")
		fileName := c.Param("filename")
		// todo validate and sanitize file name

		// Todo add role checks when auth is added
		// verify bucket exists
		var bucket models.Bucket
		if err := database.DB.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}
		// find latest object version and verify it is deleted
		object := models.Object{}
		if err := database.DB.Where("bucket_id = ? AND file_name = ? ", bucketID, fileName).
			Order("version desc").
			First(&object).Error; err != nil || !object.DeleteMarker {
			c.JSON(http.StatusNotFound, gin.H{"error": "object not found"})
			return
		}
		if err := database.DB.Unscoped().Delete(&models.Object{}, object.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "object not found"})
			return
		}
		c.IndentedJSON(http.StatusOK, object)
	}
}
