package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mic615/chill-crate-api/internal/handlers"
	"github.com/mic615/chill-crate-api/internal/models"
	"github.com/mic615/chill-crate-api/internal/testutil"
)

func TestCreateGroup(t *testing.T) {
	db, err := testutil.SetupTestDB(t)
	require.NoError(t, err)
	h := handlers.NewHandler(db, nil)
	user, err := testutil.CreateUser(db, &models.User{
		KCUserID:  "test-kc-id-123",
		FirstName: "Test",
		LastName:  "User",
		Username:  "testuser",
		Email:     "test@test.com",
	})
	require.NoError(t, err)
	body := `{"name":"Test Group","description":"This is a test group"}`
	w, c := testutil.AuthenticateRequest(&user)
	c.Request = httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/groups",
		strings.NewReader(body),
	)
	c.Request.Header.Set("Content-Type", "application/json")
	h.CreateGroup()(c)
	require.Equal(t, http.StatusCreated, w.Code)
	// verify the group was created AND the user got an admin membership
	var m models.Membership
	err = db.Where("user_id = ?", user.ID).First(&m).Error
	require.NoError(t, err)
	require.Equal(t, models.RoleAdmin, m.Role)
}

func TestGetMyGroups(t *testing.T) {
	db, err := testutil.SetupTestDB(t)
	require.NoError(t, err)
	h := handlers.NewHandler(db, nil)

	user, group := testutil.SetupUserWithGroup(t, db, models.RoleAdmin)
	require.NoError(t, err)
	w, c := testutil.AuthenticateRequest(&user)
	c.Request = httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/groups",
		http.NoBody,
	)
	h.GetMyGroups()(c)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), group.Name)
}
