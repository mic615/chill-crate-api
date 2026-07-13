package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Buckets     []Bucket       `gorm:"foreignKey:GroupID"`
	Memberships []Membership   `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE"`
}
