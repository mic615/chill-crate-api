package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mic615/chill-crate-api/internal/models"
)

func (h *Handler) UploadObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		bucketID := c.Param("bucketId")
		fileName := c.Param("filename")
		// todo validate and sanitize file name

		// Todo add role checks when auth is added
		// verify bucket exists
		var bucket models.Bucket
		if err := h.db.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}
		// calculate version
		objects := []models.Object{}
		h.db.Where("bucket_id = ? AND file_name = ? ", bucketID, fileName).
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

		if err := h.storageClient.UploadObject(
			bucket.ID.String(),
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
		if err := h.db.Create(&newObject).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, newObject)
	}
}

func (h *Handler) DownloadObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		bucketID := c.Param("bucketId")
		var bucket models.Bucket
		if err := h.db.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}
		object := models.Object{}
		// get the latest object
		if err := h.db.Where("bucket_id = ? AND file_name = ?", bucketID, filename).
			Order("version desc").
			First(&object).
			Error; err != nil || object.DeleteMarker {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "object not found"})
			return
		}
		// todo check file size
		// based on size , do single or multi part downlod
		body, err := h.storageClient.DownloadObject(bucket.ID.String(), object.StoragePath.String())
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

func (h *Handler) GetObjectsByBucketID() gin.HandlerFunc {
	return func(c *gin.Context) {
		bucketID := c.Param("bucketId")
		objects := []models.Object{}

		var bucket models.Bucket
		if err := h.db.First(&bucket, "id = ?", bucketID).Error; err != nil {
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
		if err := h.db.Raw(queryString, bucketID).Scan(&objects).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, objects)
	}
}

func (h *Handler) GetObjectByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		object := models.Object{}

		if err := h.db.First(&object, id).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, object)
	}
}

func (h *Handler) DeleteObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		bucketID := c.Param("bucketId")
		fileName := c.Param("filename")
		// todo validate and sanitize file name

		// Todo add role checks when auth is added
		// verify bucket exists
		var bucket models.Bucket
		if err := h.db.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}
		// find latest object version and verify it is not deleted already
		object := models.Object{}
		if err := h.db.Where("bucket_id = ? AND file_name = ? ", bucketID, fileName).
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
		if err := h.db.Create(&newObject).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, newObject)
	}
}

func (h *Handler) RestoreObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		bucketID := c.Param("bucketId")
		fileName := c.Param("filename")
		// todo validate and sanitize file name

		// Todo add role checks when auth is added
		// verify bucket exists
		var bucket models.Bucket
		if err := h.db.First(&bucket, "id = ?", bucketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "bucket not found"})
			return
		}
		// find last 2 objects  should be the deleted one and version n-1 (the one we're restoring)
		objects := []models.Object{}
		err := h.db.Where("bucket_id = ? AND file_name = ? ", bucketID, fileName).
			Order("version desc").Limit(2).
			Find(&objects).Error
		if len(objects) < 2 {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": "there are less than 2 object versions nothing to restore"})
			return
		}
		deletedVersion := objects[0]
		restoredVersion := objects[1]
		// find latest object version and verify it is deleted
		if err != nil || !deletedVersion.DeleteMarker {
			c.JSON(http.StatusNotFound, gin.H{"error": "object not found"})
			return
		}
		//  create a new version with the restored data
		restoredObject := models.Object{
			FileName:    restoredVersion.FileName,
			BucketID:    restoredVersion.BucketID,
			Version:     deletedVersion.Version + 1,
			StoragePath: restoredVersion.StoragePath,
			Size:        restoredVersion.Size,
		}
		//  soft delete the deleted version for auditability
		err = h.db.Transaction(func(tx *gorm.DB) error {
			if txErr := tx.Delete(&models.Object{}, deletedVersion.ID).Error; txErr != nil {
				return txErr
			}
			return tx.Create(&restoredObject).Error
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": "failed to restore object: " + err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, restoredObject)
	}
}
