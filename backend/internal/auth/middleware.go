package auth

import (
	"net/http"
	"strings"

	"fnchatbot/internal/models"

	"github.com/gin-gonic/gin"
)

const currentUserKey = "currentUser"

// getTokenFromHeader extracts Bearer token from Authorization header.
func getTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}

// JWTMiddleware authenticates HTTP requests using JWT bearer tokens.
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := getTokenFromHeader(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		var user models.User
		if err := dbConn.Where("id = ? AND enabled = ?", claims.UserID, true).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found or disabled"})
			return
		}

		c.Set(currentUserKey, &user)
		c.Next()
	}
}

// WebSocketAuthMiddleware authenticates WebSocket upgrade requests using JWT from query or header.
func WebSocketAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = getTokenFromHeader(c)
		}
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		var user models.User
		if err := dbConn.Where("id = ? AND enabled = ?", claims.UserID, true).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found or disabled"})
			return
		}

		c.Set(currentUserKey, &user)
		c.Next()
	}
}

// CasbinMiddleware performs RBAC authorization using Casbin.
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := CurrentUser(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if casbinEnf == nil {
			// If enforcer is not ready, deny by default.
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "authorization not configured"})
			return
		}

		sub := "role_user"
		if user.IsAdmin {
			sub = "role_admin"
		}

		obj := c.FullPath()
		if obj == "" {
			obj = c.Request.URL.Path
		}
		act := c.Request.Method

		allowed, err := casbinEnf.Enforce(sub, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "authorization error"})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}

// CurrentUser returns current authenticated user from context.
func CurrentUser(c *gin.Context) (*models.User, bool) {
	val, ok := c.Get(currentUserKey)
	if !ok {
		return nil, false
	}
	user, ok := val.(*models.User)
	return user, ok
}

// IsAdmin returns true if user is admin.
func IsAdmin(u *models.User) bool {
	return u != nil && u.IsAdmin
}

// IsInitialAdmin returns true if user is the initial admin (ID == 1 and admin).
func IsInitialAdmin(u *models.User) bool {
	return u != nil && u.IsAdmin && u.ID == 1
}
