package testutil

import (
	"gorm.io/gorm"

	"github.com/mic615/chill-crate-api/internal/models"
)

func CreateUser(db *gorm.DB, user *models.User) (models.User, error) {
	err := db.Create(user).Error
	return *user, err
}

func CreateGroup(db *gorm.DB, group *models.Group) (models.Group, error) {
	err := db.Create(group).Error
	return *group, err
}

func AddMembership(
	db *gorm.DB,
	u models.User,
	g models.Group,
	role models.Role,
) (models.Membership, error) {

	membership := models.Membership{
		UserID:  u.ID,
		GroupID: g.ID,
		Role:    role,
	}
	err := db.Create(&membership).Error
	return membership, err
}

func CreateBucket(db *gorm.DB, bucket *models.Bucket) (models.Bucket, error) {
	err := db.Create(bucket).Error
	return *bucket, err
}
