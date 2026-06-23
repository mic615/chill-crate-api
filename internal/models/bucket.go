package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bucket struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null;uniqueIndex:idx_bucket_name_group"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_bucket_name_group"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Objects   []Object
}

func (b *Bucket) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
