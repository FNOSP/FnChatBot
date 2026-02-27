package auth

import (
	"fmt"
	"time"

	"fnchatbot/internal/config"
	"fnchatbot/internal/models"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Claims represents JWT claims for authenticated users.
type Claims struct {
	UserID             uint   `json:"user_id"`
	Username           string `json:"username"`
	IsAdmin            bool   `json:"is_admin"`
	MustChangePassword bool   `json:"must_change_password"`
	jwt.RegisteredClaims
}

var (
	dbConn    *gorm.DB
	jwtSecret []byte
	tokenTTL  time.Duration
	casbinEnf *casbin.Enforcer
)

// Init initializes authentication (JWT, admin user) and Casbin authorization.
func Init(db *gorm.DB, authCfg config.AuthConfig) error {
	dbConn = db
	if authCfg.JWTSecret == "" {
		return fmt.Errorf("auth.jwt_secret must be configured")
	}
	jwtSecret = []byte(authCfg.JWTSecret)
	if authCfg.TokenLifetime == 0 {
		tokenTTL = 24 * time.Hour
	} else {
		tokenTTL = time.Duration(authCfg.TokenLifetime) * time.Second
	}

	// Initialize initial admin user from config.
	if err := InitAdminUser(db, authCfg.InitialAdmin); err != nil {
		return err
	}

	// Initialize Casbin enforcer.
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	enf, err := casbin.NewEnforcer("internal/auth/casbin_model.conf", adapter)
	if err != nil {
		return fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// Load existing policies from DB.
	if err := enf.LoadPolicy(); err != nil {
		return fmt.Errorf("failed to load casbin policy: %w", err)
	}

	// Seed basic policies (idempotent in Casbin).
	_, _ = enf.AddPolicy("role_admin", "/api/*", "(GET|POST|PUT|PATCH|DELETE)")
	_, _ = enf.AddPolicy("role_admin", "/ws/chat/*", "GET")

	// Normal user can use most read/write APIs but not user/sandbox admin endpoints.
	_, _ = enf.AddPolicy("role_user", "/api/auth/*", "(GET|POST)")
	_, _ = enf.AddPolicy("role_user", "/api/conversations*", "(GET|POST|DELETE)")
	_, _ = enf.AddPolicy("role_user", "/api/models*", "(GET|POST)")
	_, _ = enf.AddPolicy("role_user", "/api/skills*", "(GET|POST|PATCH|DELETE)")
	_, _ = enf.AddPolicy("role_user", "/api/mcp*", "(GET|POST|PUT|DELETE)")
	_, _ = enf.AddPolicy("role_user", "/api/providers*", "(GET|POST|PUT)")
	_, _ = enf.AddPolicy("role_user", "/api/sandbox", "GET")
	_, _ = enf.AddPolicy("role_user", "/api/sandbox/paths*", "GET")
	if err := enf.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save casbin policy: %w", err)
	}

	casbinEnf = enf
	return nil
}

// InitAdminUser seeds the initial admin account if no admin exists.
func InitAdminUser(db *gorm.DB, cfg config.InitialAdminConfig) error {
	// If config is not set, do nothing.
	if cfg.Username == "" || cfg.Password == "" {
		return nil
	}

	var count int64
	if err := db.Model(&models.User{}).Where("is_admin = ?", true).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash initial admin password: %w", err)
	}

	user := models.User{
		Username:           cfg.Username,
		PasswordHash:       string(hash),
		IsAdmin:            true,
		Enabled:            true,
		MustChangePassword: true,
	}
	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create initial admin user: %w", err)
	}

	return nil
}

// HashPassword hashes plain text password using bcrypt.
func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares hashed password with plain text.
func CheckPassword(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
