package auth

import (
	"errors"
	"time"

	"fnchatbot/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken generates a JWT for the given user.
func GenerateToken(user *models.User) (string, error) {
	if jwtSecret == nil {
		return "", errors.New("jwt secret not initialized")
	}

	now := time.Now()
	claims := &Claims{
		UserID:             user.ID,
		Username:           user.Username,
		IsAdmin:            user.IsAdmin,
		MustChangePassword: user.MustChangePassword,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken parses and validates a JWT string.
func ParseToken(tokenString string) (*Claims, error) {
	if jwtSecret == nil {
		return nil, errors.New("jwt secret not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
