package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bucket struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null;uniqueIndex:idx_bucket_name_group"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Objects   []Object
}
