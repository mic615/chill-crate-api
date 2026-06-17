package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Buckets     []Bucket
	Memberships []Membership
}

func (g *Group) BeforeCreate(tx *gorm.DB) error {
	g.ID = uuid.New()
	return nil
}
