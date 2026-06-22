package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Object struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	FileName     string     `gorm:"not null;uniqueIndex:idx_object_filename_bucket_version"`
	BucketID     uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_object_filename_bucket_version"`
	Version      int        `gorm:"not null;default:1;uniqueIndex:idx_object_filename_bucket_version"`
	StoragePath  *uuid.UUID `gorm:"type:uuid"`
	Size         int64      `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	DeleteMarker bool           `gorm:"not null;default:false"`
}

func (o *Object) BeforeCreate(tx *gorm.DB) error {
	o.ID = uuid.New()
	return nil
}
