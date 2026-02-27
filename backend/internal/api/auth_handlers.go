package api

import (
	"net/http"

	"fnchatbot/internal/auth"
	"fnchatbot/internal/db"
	"fnchatbot/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPublicAuthRoutes registers auth routes that require no authentication.
func RegisterPublicAuthRoutes(r *gin.RouterGroup) {
	r.POST("/auth/login", loginHandler)
}

// RegisterAuthProtectedRoutes registers auth routes that need JWT but not Casbin.
func RegisterAuthProtectedRoutes(r *gin.RouterGroup) {
	r.GET("/auth/me", meHandler)
	r.POST("/auth/reset-password", resetPasswordHandler)
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token              string       `json:"token"`
	User               *models.User `json:"user"`
	MustChangePassword bool         `json:"must_change_password"`
}

func loginHandler(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := getDB(c).Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !user.Enabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "user disabled"})
		return
	}

	if err := auth.CheckPassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Token:              token,
		User:               &user,
		MustChangePassword: user.MustChangePassword,
	})
}

func meHandler(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, user)
}

type resetPasswordRequest struct {
	TargetUserID       uint   `json:"target_user_id"`
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password" binding:"required"`
	NewPasswordConfirm string `json:"new_password_confirm" binding:"required"`
}

func resetPasswordHandler(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NewPassword != req.NewPasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	db := getDB(c)

	targetID := req.TargetUserID
	if targetID == 0 {
		targetID = user.ID
	}

	var target models.User
	if err := db.First(&target, targetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Permission checks:
	switch {
	case target.ID == user.ID:
		// Self password change.
		if !auth.IsInitialAdmin(user) {
			// Non-initial users must provide old password.
			if err := auth.CheckPassword(user.PasswordHash, req.OldPassword); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
				return
			}
		}
	case auth.IsAdmin(user) && !target.IsAdmin:
		// Admin changing normal user's password: allowed without old password.
	case auth.IsInitialAdmin(user) && target.IsAdmin:
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

	updates := map[string]interface{}{
		"password_hash":        hash,
		"must_change_password": false,
	}

	if err := db.Model(&models.User{}).Where("id = ?", target.ID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

func getDB(c *gin.Context) *gorm.DB {
	return db.DB
}
