package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mic615/chill-crate-api/internal/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", handlers.Ping())
	groups := r.Group("/groups")
	groups.POST("", handlers.CreateGroup())
	groups.GET("", handlers.GetMyGroups())
	groups.GET("/:groupId/buckets", handlers.GetBucketsByGroupID())

	// buckets
	buckets := r.Group("/buckets")
	buckets.POST("", handlers.CreateBucket())
	buckets.GET("/:bucketId/objects", handlers.GetObjectsByBucketID())
	buckets.GET("/:bucketId/objects/:filename", handlers.DownloadObject())
	buckets.POST("/:bucketId/objects", handlers.UploadObject())

	// objects
	objects := r.Group("/objects")
	objects.GET("/:id", handlers.GetObjectByID())
}
