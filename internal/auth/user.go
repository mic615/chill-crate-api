package auth

import (
	"errors"

	"gorm.io/gorm"

	"github.com/mic615/chill-crate-api/internal/models"
)

func (a *Authenticator) resolveOrCreateUser(claims *TokenClaims) (*models.User, error) {
	var user models.User
	// Check if the user already exists in the database
	if err := a.db.Where("kc_user_id = ?", claims.Sub).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User does not exist, create a new user
			user = models.User{
				KCUserID:  claims.Sub,
				Email:     claims.Email,
				FirstName: claims.FirstName,
				LastName:  claims.LastName,
				Username:  claims.Username,
			}
			if createErr := a.db.Create(&user).Error; createErr != nil {
				return nil, createErr
			}
		} else {
			return nil, err
		}
	}

	return &user, nil
}
