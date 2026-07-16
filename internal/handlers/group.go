package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/mic615/chill-crate-api/internal/models"
)

type NewGroup struct {
	Name string `json:"name" binding:"required"`
}

type NewMember struct {
	Identifier string      `json:"identifier" binding:"required"` // email or username
	Role       models.Role `json:"role"       binding:"required"`
}

var ErrLastAdmin = errors.New("can't demote the last admin")

func (h *Handler) CreateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var group NewGroup
		user, exists := c.Get("user")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		userID := user.(*models.User).ID
		if err := c.ShouldBindJSON(&group); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newGroup := models.Group{Name: group.Name}
		err := h.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&newGroup).Error; err != nil {
				return err
			}
			membership := models.Membership{
				UserID:  userID,
				GroupID: newGroup.ID,
				Role:    models.RoleAdmin,
			}
			return tx.Create(&membership).Error
		})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, newGroup)
	}
}

func (h *Handler) GetMyGroups() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		userID := user.(*models.User).ID
		groups := []models.Group{}
		if err := h.db.Joins("JOIN memberships ON memberships.group_id = groups.id").
			Where("memberships.user_id = ?", userID).
			Find(&groups).
			Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, groups)
	}
}

func (h *Handler) GetGroups() gin.HandlerFunc {
	return func(c *gin.Context) {
		groups := []models.Group{}
		if err := h.db.Find(&groups).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, groups)
	}
}

func (h *Handler) AddMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		var member NewMember
		if err := c.ShouldBindJSON(&member); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		groupID := c.Param("groupId")
		var group models.Group
		if err := h.db.First(&group, "id = ?", groupID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		// RBAC
		if !h.authorize(c, group.ID, models.RoleAdmin) {
			return
		}
		newUser, err := h.findUserByIdentifier(member.Identifier)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		newMember := models.Membership{UserID: newUser.ID, GroupID: group.ID, Role: member.Role}
		if err := h.db.Create(&newMember).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				c.IndentedJSON(http.StatusConflict, gin.H{"error": "user already a member"})
				return
			}
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, newMember)
	}
}

func (h *Handler) UpdateRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var member NewMember
		if err := c.ShouldBindJSON(&member); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		groupID := c.Param("groupId")
		var group models.Group
		if err := h.db.First(&group, "id = ?", groupID).Error; err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		// RBAC
		if !h.authorize(c, group.ID, models.RoleAdmin) {
			return
		}
		newUser, err := h.findUserByIdentifier(member.Identifier)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		membership, err := h.updateMembershipRole(group.ID, newUser.ID, member.Role)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "membership not found"})
			case errors.Is(err, ErrLastAdmin):
				c.IndentedJSON(http.StatusConflict, gin.H{"error": err.Error()})
			default:
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.IndentedJSON(http.StatusOK, membership)
	}
}

// findUserByIdentifier looks up a user by email or by username otherwise.
func (h *Handler) findUserByIdentifier(identifier string) (models.User, error) {
	var user models.User
	query := "username = ?"
	if strings.Contains(identifier, "@") {
		query = "email = ?"
	}
	err := h.db.First(&user, query, identifier).Error
	return user, err
}

// updateMembershipRole changes a user's role within a group, refusing to
// demote the group's last remaining admin.
func (h *Handler) updateMembershipRole(
	groupID, userID uuid.UUID,
	newRole models.Role,
) (models.Membership, error) {
	var membership models.Membership
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// lock the membership row for this group to serialize concurrent demotions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND group_id = ?", userID, groupID).
			First(&membership).Error; err != nil {
			return err
		}
		// if this change removes an admin, ensure at least one admin remains
		if membership.Role == models.RoleAdmin && newRole != models.RoleAdmin {
			var adminCount int64
			err := tx.Model(&models.Membership{}).
				Where("group_id = ? AND role = ?", groupID, models.RoleAdmin).
				Count(&adminCount).Error
			if err != nil {
				return err
			}
			if adminCount <= 1 {
				return ErrLastAdmin
			}
		}
		membership.Role = newRole
		return tx.Save(&membership).Error
	})
	return membership, err
}
