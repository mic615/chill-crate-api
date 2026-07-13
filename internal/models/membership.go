package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// todo add roles
type Membership struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_membership_user_group"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_membership_user_group"`
	Role      Role      `gorm:"not null;default:viewer"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User  `gorm:"foreignKey:UserID"`
	Group     Group `gorm:"foreignKey:GroupID"`
}

func (m *Membership) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	return nil
}
