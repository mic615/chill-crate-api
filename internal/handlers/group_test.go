package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

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

func newMemberUser(t *testing.T, db *gorm.DB) models.User {
	t.Helper()
	user, err := testutil.CreateUser(db, &models.User{
		KCUserID:  "kc-" + uuid.NewString(),
		Email:     uuid.NewString() + "@example.com",
		FirstName: "New",
		LastName:  "Member",
		Username:  uuid.NewString(),
	})
	require.NoError(t, err)
	return user
}

func addMemberRequest(
	admin *models.User,
	groupID, body string,
) (*httptest.ResponseRecorder, *gin.Context) {
	w, c := testutil.AuthenticateRequest(admin)
	c.Request = httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/groups/"+groupID+"/members",
		strings.NewReader(body),
	)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "groupId", Value: groupID}}
	return w, c
}

func TestAddMember(t *testing.T) {
	db, err := testutil.SetupTestDB(t)
	require.NoError(t, err)
	h := handlers.NewHandler(db, nil)

	t.Run("admin adds member by email", func(t *testing.T) {
		admin, group := testutil.SetupUserWithGroup(t, db, models.RoleAdmin)
		newUser := newMemberUser(t, db)

		body := `{"identifier":"` + newUser.Email + `","role":"editor"}`
		w, c := addMemberRequest(&admin, group.ID.String(), body)
		h.AddMember()(c)
		require.Equal(t, http.StatusCreated, w.Code)

		var m models.Membership
		require.NoError(
			t,
			db.Where("user_id = ? AND group_id = ?", newUser.ID, group.ID).First(&m).Error,
		)
		require.Equal(t, models.RoleEditor, m.Role)
	})

	t.Run("admin adds member by username", func(t *testing.T) {
		admin, group := testutil.SetupUserWithGroup(t, db, models.RoleAdmin)
		newUser := newMemberUser(t, db)

		body := `{"identifier":"` + newUser.Username + `","role":"editor"}`
		w, c := addMemberRequest(&admin, group.ID.String(), body)
		h.AddMember()(c)
		require.Equal(t, http.StatusCreated, w.Code)

		var m models.Membership
		require.NoError(
			t,
			db.Where("user_id = ? AND group_id = ?", newUser.ID, group.ID).First(&m).Error,
		)
		require.Equal(t, models.RoleEditor, m.Role)
	})

	t.Run("non-admin forbidden", func(t *testing.T) {
		editor, group := testutil.SetupUserWithGroup(t, db, models.RoleEditor)
		newUser := newMemberUser(t, db)

		body := `{"identifier":"` + newUser.Email + `","role":"editor"}`
		w, c := addMemberRequest(&editor, group.ID.String(), body)
		h.AddMember()(c)
		require.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("group not found", func(t *testing.T) {
		admin, _ := testutil.SetupUserWithGroup(t, db, models.RoleAdmin)
		newUser := newMemberUser(t, db)

		body := `{"identifier":"` + newUser.Email + `","role":"editor"}`
		w, c := addMemberRequest(&admin, uuid.NewString(), body)
		h.AddMember()(c)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("target user not found", func(t *testing.T) {
		admin, group := testutil.SetupUserWithGroup(t, db, models.RoleAdmin)

		body := `{"identifier":"nobody@example.com","role":"editor"}`
		w, c := addMemberRequest(&admin, group.ID.String(), body)
		h.AddMember()(c)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("duplicate member conflict", func(t *testing.T) {
		admin, group := testutil.SetupUserWithGroup(t, db, models.RoleAdmin)
		existingUser := newMemberUser(t, db)
		_, err := testutil.AddMembership(db, existingUser, group, models.RoleEditor)
		require.NoError(t, err)

		body := `{"identifier":"` + existingUser.Email + `","role":"editor"}`
		w, c := addMemberRequest(&admin, group.ID.String(), body)
		h.AddMember()(c)
		require.Equal(t, http.StatusConflict, w.Code)
	})
}
