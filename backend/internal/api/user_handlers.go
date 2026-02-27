package api

import (
	"net/http"
	"strconv"

	"fnchatbot/internal/auth"
	"fnchatbot/internal/db"
	"fnchatbot/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUsers returns users according to role: admin sees all, normal user sees only self.
func GetUsers(c *gin.Context) {
	current, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var users []models.User
	if auth.IsAdmin(current) {
		if err := db.DB.Order("id asc").Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := db.DB.Where("id = ?", current.ID).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, users)
}

// GetCurrentUser returns current user profile.
func GetCurrentUser(c *gin.Context) {
	current, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, current)
}

type createUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type"`    // "admin" or "user"
	Enabled     *bool  `json:"enabled"` // optional
}

// CreateUser creates a new user (admin only).
func CreateUser(c *gin.Context) {
	current, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !auth.IsAdmin(current) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var enabled = true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	isAdmin := req.Type == "admin"

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{
		Username:     req.Username,
		PasswordHash: hash,
		Description:  req.Description,
		IsAdmin:      isAdmin,
		Enabled:      enabled,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

type updateUserRequest struct {
	Username    *string `json:"username"`
	Description *string `json:"description"`
	Enabled     *bool   `json:"enabled"`
}

// UpdateUser updates user basic fields with role checks.
func UpdateUser(c *gin.Context) {
	current, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var target models.User
	if err := db.DB.First(&target, targetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Permission rules:
	// - Normal user: only update own description.
	// - Admin:
	//   - On normal users: can edit username/description/enabled.
	//   - On admin users:
	//     - Initial admin can edit other admins.
	//     - Other admins can only edit self.
	updates := map[string]interface{}{}

	if !auth.IsAdmin(current) {
		// Normal user can only update own description.
		if current.ID != target.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if req.Description != nil {
			updates["description"] = *req.Description
		}
	} else {
		// Admin path.
		if target.IsAdmin && !auth.IsInitialAdmin(current) && current.ID != target.ID {
			// Non-initial admin cannot edit other admins.
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		if req.Username != nil {
			updates["username"] = *req.Username
		}
		if req.Description != nil {
			updates["description"] = *req.Description
		}
		if req.Enabled != nil {
			// Prevent disabling the last admin.
			if target.IsAdmin && !*req.Enabled {
				var adminCount int64
				if err := db.DB.Model(&models.User{}).Where("is_admin = ? AND enabled = ?", true, true).Count(&adminCount).Error; err == nil && adminCount <= 1 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "cannot disable the last admin"})
					return
				}
			}
			updates["enabled"] = *req.Enabled
		}
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, target)
		return
	}

	if err := db.DB.Model(&models.User{}).Where("id = ?", target.ID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.First(&target, target.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, target)
}

type changePasswordRequest struct {
	NewPassword        string `json:"new_password" binding:"required"`
	NewPasswordConfirm string `json:"new_password_confirm" binding:"required"`
	OldPassword        string `json:"old_password"`
}

// ChangeUserPassword allows changing user passwords based on role rules.
func ChangeUserPassword(c *gin.Context) {
	current, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	targetID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var target models.User
	if err := db.DB.First(&target, targetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NewPassword != req.NewPasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	switch {
	case target.ID == current.ID:
		// Self change.
		if !auth.IsInitialAdmin(current) {
			// Require old password for non-initial users.
			if err := auth.CheckPassword(current.PasswordHash, req.OldPassword); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
				return
			}
		}
	case auth.IsAdmin(current) && !target.IsAdmin:
		// Admin changing normal user password: allowed.
	case auth.IsInitialAdmin(current) && target.IsAdmin:
		// Initial admin changing other admin's password.
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	if err := db.DB.Model(&models.User{}).Where("id = ?", target.ID).Update("password_hash", hash).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}
