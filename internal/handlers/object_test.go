package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/mic615/chill-crate-api/internal/models"
	"github.com/mic615/chill-crate-api/internal/testutil"
)

func TestUploadObject(t *testing.T) {
	db, err := testutil.SetupTestDB(t)
	storage := testutil.SetupTestStorageClient(t)
	require.NoError(t, err)
	h := NewHandler(db, storage)

	user, group := testutil.SetupUserWithGroup(t, db, models.RoleAdmin)
	bucket, err := testutil.CreateBucket(db, &models.Bucket{
		Name:    "Test Bucket",
		GroupID: group.ID,
	})
	require.NoError(t, err)

	body := `{ "content": "This is a test object" }`

	w, c := testutil.AuthenticateRequest(&user)
	c.Request = httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/buckets/"+bucket.ID.String()+"/objects/test-object.txt",
		strings.NewReader(body),
	)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "bucketId", Value: bucket.ID.String()},
		{Key: "filename", Value: "test-object.txt"},
	}
	h.UploadObject()(c)
	require.Equal(t, http.StatusCreated, w.Code)
	// verify the object was created
	var object models.Object
	err = db.Where("bucket_id = ?", bucket.ID).First(&object).Error
	require.NoError(t, err)
	require.Equal(t, "test-object.txt", object.FileName)
}
