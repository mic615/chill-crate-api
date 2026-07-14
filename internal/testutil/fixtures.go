package testutil

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/mic615/chill-crate-api/internal/models"
)

func SetupUserWithGroup(t *testing.T, db *gorm.DB, role models.Role) (models.User, models.Group) {
	t.Helper()
	user, err := CreateUser(
		db,
		&models.User{
			KCUserID:  "kc-" + uuid.NewString(),
			Email:     uuid.NewString() + "@example.com",
			FirstName: "Test",
			LastName:  "User",
			Username:  uuid.NewString(),
		},
	)
	require.NoError(t, err)
	group, err := CreateGroup(db, &models.Group{Name: "Test Group"})
	require.NoError(t, err)
	_, err = AddMembership(db, user, group, role)
	require.NoError(t, err)
	return user, group
}
