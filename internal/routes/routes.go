package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mic615/chill-crate-api/internal/auth"
	"github.com/mic615/chill-crate-api/internal/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", handlers.Ping())

	// everything below requires auth
	protected := r.Group("/api")
	protected.Use(auth.AuthMiddleware())
	// groups
	groups := protected.Group("/groups")
	groups.POST("", handlers.CreateGroup())
	groups.GET("", handlers.GetMyGroups())
	groups.GET("/:groupId/buckets", handlers.GetBucketsByGroupID())
	groups.GET("/:groupId/buckets/:name", handlers.GetBucketByName())

	// buckets
	buckets := protected.Group("/buckets")
	buckets.POST("", handlers.CreateBucket())
	buckets.GET("/:bucketId/objects", handlers.GetObjectsByBucketID())
	buckets.GET("/:bucketId/objects/:filename", handlers.DownloadObject())
	buckets.POST("/:bucketId/objects/:filename", handlers.UploadObject())
	buckets.DELETE("/:bucketId/objects/:filename", handlers.DeleteObject())
	buckets.POST("/:bucketId/objects/:filename/restore", handlers.RestoreObject())
	buckets.DELETE("/:bucketId", handlers.DeleteBucket())

	// objects
	objects := protected.Group("/objects")
	objects.GET("/:id", handlers.GetObjectByID())
}
