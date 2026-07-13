package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	KCUserID    string    `gorm:"not null;unique"`
	Username    string    `gorm:"not null;unique"`
	Email       string    `gorm:"not null;unique"`
	FirstName   string
	LastName    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Memberships []Membership   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
