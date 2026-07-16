package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/mic615/chill-crate-api/internal/handlers"
	"github.com/mic615/chill-crate-api/internal/models"
	"github.com/mic615/chill-crate-api/internal/testutil"
)

func TestCreateBucket(t *testing.T) {
	tests := []struct {
		name     string
		role     models.Role
		wantCode int
	}{
		{"editor allowed", models.RoleEditor, http.StatusCreated},
		{"admin allowed", models.RoleAdmin, http.StatusCreated},
		{"viewer forbidden", models.RoleViewer, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := testutil.SetupTestDB(t)
			require.NoError(t, err)
			storage := testutil.SetupTestStorageClient(t)
			h := handlers.NewHandler(db, storage)
			user, group := testutil.SetupUserWithGroup(t, db, tt.role)

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
			require.Equal(t, tt.wantCode, w.Code)

			var bucket models.Bucket
			err = db.Where("group_id = ?", group.ID).First(&bucket).Error
			if tt.wantCode == http.StatusCreated {
				require.NoError(t, err)
				require.Equal(t, "Test Bucket", bucket.Name)
			} else {
				require.ErrorIs(t, err, gorm.ErrRecordNotFound)
			}
		})
	}
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
