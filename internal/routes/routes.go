package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mic615/chill-crate-api/internal/handlers"
)

func RegisterRoutes(r *gin.Engine, h *handlers.Handler, authMiddleware gin.HandlerFunc) {
	r.GET("/ping", handlers.Ping())

	// everything below requires auth
	protected := r.Group("/api")
	if authMiddleware == nil {
		authMiddleware = func(c *gin.Context) { c.Next() }
	}
	protected.Use(authMiddleware)
	// groups
	groups := protected.Group("/groups")
	groups.POST("", h.CreateGroup())
	groups.GET("", h.GetMyGroups())
	groups.POST("/:groupId/members", h.AddMember())
	groups.PUT("/:groupId/members", h.UpdateRole())
	groups.GET("/:groupId/buckets", h.GetBucketsByGroupID())
	groups.GET("/:groupId/buckets/:name", h.GetBucketByName())

	// buckets
	buckets := protected.Group("/buckets")
	buckets.POST("", h.CreateBucket())
	buckets.GET("/:bucketId/objects", h.GetObjectsByBucketID())
	buckets.GET("/:bucketId/objects/:filename", h.DownloadObject())
	buckets.POST("/:bucketId/objects/:filename", h.UploadObject())
	buckets.DELETE("/:bucketId/objects/:filename", h.DeleteObject())
	buckets.POST("/:bucketId/objects/:filename/restore", h.RestoreObject())
	buckets.DELETE("/:bucketId", h.DeleteBucket())

	// objects
	objects := protected.Group("/objects")
	objects.GET("/:id", h.GetObjectByID())
}
