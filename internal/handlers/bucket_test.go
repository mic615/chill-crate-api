package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/mic615/chill-crate-api/internal/handlers"
	"github.com/mic615/chill-crate-api/internal/models"
	"github.com/mic615/chill-crate-api/internal/testutil"
)

func TestCreateBucket(t *testing.T) {
	db, err := testutil.SetupTestDB(t)
	storage := testutil.SetupTestStorageClient(t)
	require.NoError(t, err)
	h := handlers.NewHandler(db, storage)
	user, group := testutil.SetupUserWithGroup(t, db, models.RoleAdmin)

	body := `{"name":"Test Bucket","group_id":"` + group.ID.String() + `"}`
	w, c := testutil.AuthenticateRequest(&user)
	c.Request = httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/buckets",
		strings.NewReader(body),
	)
	c.Request.Header.Set("Content-Type", "application/json")
	h.CreateBucket()(c)
	require.Equal(t, http.StatusCreated, w.Code)
	// verify the bucket was created
	var bucket models.Bucket
	err = db.Where("group_id = ?", group.ID).First(&bucket).Error
	require.NoError(t, err)
	require.Equal(t, "Test Bucket", bucket.Name)
}

func TestGetBucketsByGroupID(t *testing.T) {
	db, err := testutil.SetupTestDB(t)
	require.NoError(t, err)
	h := handlers.NewHandler(db, nil)
	user, group := testutil.SetupUserWithGroup(t, db, models.RoleAdmin)
	require.NoError(t, err)
	bucket, err := testutil.CreateBucket(db, &models.Bucket{
		Name:    "Test Bucket",
		GroupID: group.ID,
	})
	require.NoError(t, err)

	w, c := testutil.AuthenticateRequest(&user)
	c.Request = httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/groups/"+group.ID.String()+"/buckets",
		http.NoBody,
	)
	c.Params = gin.Params{{Key: "groupId", Value: group.ID.String()}}
	h.GetBucketsByGroupID()(c)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), bucket.Name)
}
