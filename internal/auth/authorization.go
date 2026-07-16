package auth

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mic615/chill-crate-api/internal/models"
)

var ErrForbidden = errors.New("unauthorized")

func RequireRole(db *gorm.DB, userID, groupID uuid.UUID, role models.Role) error {
	var m models.Membership
	err := db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrForbidden // not a member of the group
		}
		return err
	}
	if m.Role < role { // roles are least privilege viewer=0 < editor=1 < admin=2
		return ErrForbidden
	}
	return nil
}
