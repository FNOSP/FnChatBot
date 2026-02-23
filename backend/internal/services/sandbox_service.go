package services

import (
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"fnchatbot/internal/models"

	"gorm.io/gorm"
)

type SandboxService struct {
	db *gorm.DB
}

type SandboxPath struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Path        string `gorm:"type:varchar(500);not null;uniqueIndex" json:"path"`
	Description string `gorm:"type:varchar(500)" json:"description"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`
}

type SandboxConfig struct {
	ID      uint `gorm:"primaryKey" json:"id"`
	Enabled bool `gorm:"default:false" json:"enabled"`
}

func NewSandboxService(db *gorm.DB) *SandboxService {
	db.AutoMigrate(&SandboxPath{}, &SandboxConfig{})
	return &SandboxService{db: db}
}

func (s *SandboxService) IsEnabled() bool {
	var config SandboxConfig
	if err := s.db.First(&config).Error; err != nil {
		return false
	}
	return config.Enabled
}

func (s *SandboxService) GetAllowedPaths() []string {
	var paths []SandboxPath
	s.db.Where("enabled = ?", true).Find(&paths)
	result := make([]string, len(paths))
	for i, p := range paths {
		result[i] = p.Path
	}
	return result
}

func (s *SandboxService) AddPath(path string, description string) error {
	absPath, err := s.normalizePath(path)
	if err != nil {
		return err
	}
	sandboxPath := SandboxPath{
		Path:        absPath,
		Description: description,
		Enabled:     true,
	}
	return s.db.Create(&sandboxPath).Error
}

func (s *SandboxService) RemovePath(path string) error {
	absPath, err := s.normalizePath(path)
	if err != nil {
		return err
	}
	return s.db.Where("path = ?", absPath).Delete(&SandboxPath{}).Error
}

func (s *SandboxService) IsPathAllowed(path string) bool {
	if !s.IsEnabled() {
		return true
	}

	absPath, err := s.normalizePath(path)
	if err != nil {
		return false
	}

	allowedPaths := s.GetAllowedPaths()
	for _, allowed := range allowedPaths {
		if s.isSubPath(allowed, absPath) {
			return true
		}
	}
	return false
}

func (s *SandboxService) normalizePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", nil
	}

	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "/", "\\")
		if len(path) >= 2 && path[1] == ':' {
			path = strings.ToUpper(string(path[0])) + path[1:]
		}
	} else {
		path = strings.ReplaceAll(path, "\\", "/")
	}

	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return path, nil
		}
		path = absPath
	}

	return filepath.Clean(path), nil
}

func (s *SandboxService) isSubPath(parent, child string) bool {
	parent = strings.ToLower(filepath.Clean(parent))
	child = strings.ToLower(filepath.Clean(child))

	if parent == child {
		return true
	}

	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(parent, string(filepath.Separator)) {
			parent += string(filepath.Separator)
		}
		return strings.HasPrefix(child, parent)
	}

	if !strings.HasSuffix(parent, "/") {
		parent += "/"
	}
	return strings.HasPrefix(child, parent)
}

func (s *SandboxService) ExtractPathsFromCommand(command string) []string {
	paths := []string{}

	pathPatterns := []struct {
		cmd   string
		regex *regexp.Regexp
	}{
		{"cd", regexp.MustCompile(`(?i)\bcd\s+["']?([^\s"']+)["']?`)},
		{"ls", regexp.MustCompile(`(?i)\bls\s+(?:-[a-zA-Z]+\s+)*["']?([^\s"']+)["']?`)},
		{"cat", regexp.MustCompile(`(?i)\bcat\s+["']?([^\s"']+)["']?`)},
		{"rm", regexp.MustCompile(`(?i)\brm\s+(?:-[a-zA-Z]+\s+)*["']?([^\s"']+)["']?`)},
		{"cp", regexp.MustCompile(`(?i)\bcp\s+(?:-[a-zA-Z]+\s+)*["']?([^\s"']+)["']?\s+["']?([^\s"']+)["']?`)},
		{"mv", regexp.MustCompile(`(?i)\bmv\s+(?:-[a-zA-Z]+\s+)*["']?([^\s"']+)["']?\s+["']?([^\s"']+)["']?`)},
		{"mkdir", regexp.MustCompile(`(?i)\bmkdir\s+(?:-[a-zA-Z]+\s+)*["']?([^\s"']+)["']?`)},
		{"touch", regexp.MustCompile(`(?i)\btouch\s+["']?([^\s"']+)["']?`)},
		{"chmod", regexp.MustCompile(`(?i)\bchmod\s+(?:-[a-zA-Z]+\s+)?\d+\s+["']?([^\s"']+)["']?`)},
		{"chown", regexp.MustCompile(`(?i)\bchown\s+(?:-[a-zA-Z]+\s+)?[^\s]+\s+["']?([^\s"']+)["']?`)},
		{"find", regexp.MustCompile(`(?i)\bfind\s+["']?([^\s"']+)["']?`)},
		{"grep", regexp.MustCompile(`(?i)\bgrep\s+(?:-[a-zA-Z]+\s+)*["']?([^\s"']+)["']?`)},
		{"head", regexp.MustCompile(`(?i)\bhead\s+(?:-[a-zA-Z]+\s+)*["']?([^\s"']+)["']?`)},
		{"tail", regexp.MustCompile(`(?i)\btail\s+(?:-[a-zA-Z]+\s+)*["']?([^\s"']+)["']?`)},
		{"less", regexp.MustCompile(`(?i)\bless\s+["']?([^\s"']+)["']?`)},
		{"more", regexp.MustCompile(`(?i)\bmore\s+["']?([^\s"']+)["']?`)},
		{"nano", regexp.MustCompile(`(?i)\bnano\s+["']?([^\s"']+)["']?`)},
		{"vim", regexp.MustCompile(`(?i)\bvim\s+["']?([^\s"']+)["']?`)},
		{"vi", regexp.MustCompile(`(?i)\bvi\s+["']?([^\s"']+)["']?`)},
		{"echo", regexp.MustCompile(`(?i)\becho\s+.*?>>\s*["']?([^\s"']+)["']?`)},
		{"type", regexp.MustCompile(`(?i)\btype\s+["']?([^\s"']+)["']?`)},
		{"dir", regexp.MustCompile(`(?i)\bdir\s+["']?([^\s"']+)["']?`)},
		{"del", regexp.MustCompile(`(?i)\bdel\s+["']?([^\s"']+)["']?`)},
		{"copy", regexp.MustCompile(`(?i)\bcopy\s+["']?([^\s"']+)["']?`)},
		{"move", regexp.MustCompile(`(?i)\bmove\s+["']?([^\s"']+)["']?`)},
		{"xcopy", regexp.MustCompile(`(?i)\bxcopy\s+["']?([^\s"']+)["']?`)},
	}

	for _, pattern := range pathPatterns {
		matches := pattern.regex.FindAllStringSubmatch(command, -1)
		for _, match := range matches {
			for i := 1; i < len(match); i++ {
				if match[i] != "" && !s.isFlag(match[i]) {
					paths = append(paths, match[i])
				}
			}
		}
	}

	genericPathRegex := regexp.MustCompile(`["']([A-Za-z]:[\\/][^"']+|/(?:[^/"']+/)*[^"']+)["']`)
	matches := genericPathRegex.FindAllStringSubmatch(command, -1)
	for _, match := range matches {
		if len(match) > 1 && match[1] != "" {
			paths = append(paths, match[1])
		}
	}

	unquotedPathRegex := regexp.MustCompile(`\s([A-Za-z]:[\\/][^\s]+|/(?:[^/\s]+/)+[^\s]*)`)
	matches = unquotedPathRegex.FindAllStringSubmatch(command, -1)
	for _, match := range matches {
		if len(match) > 1 && match[1] != "" && !s.isFlag(match[1]) {
			paths = append(paths, match[1])
		}
	}

	return s.uniquePaths(paths)
}

func (s *SandboxService) isFlag(str string) bool {
	return strings.HasPrefix(str, "-") || strings.HasPrefix(str, "/")
}

func (s *SandboxService) uniquePaths(paths []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, path := range paths {
		normalized, err := s.normalizePath(path)
		if err != nil {
			normalized = path
		}
		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, path)
		}
	}
	return result
}

func (s *SandboxService) CheckCommandPermission(command string) (allowed bool, blockedPaths []string) {
	if !s.IsEnabled() {
		return true, nil
	}

	paths := s.ExtractPathsFromCommand(command)
	if len(paths) == 0 {
		return true, nil
	}

	for _, path := range paths {
		if !s.IsPathAllowed(path) {
			blockedPaths = append(blockedPaths, path)
		}
	}

	if len(blockedPaths) > 0 {
		return false, blockedPaths
	}

	return true, nil
}

func (s *SandboxService) SetEnabled(enabled bool) error {
	var config SandboxConfig
	result := s.db.First(&config)
	if result.Error == gorm.ErrRecordNotFound {
		config = SandboxConfig{Enabled: enabled}
		return s.db.Create(&config).Error
	}
	return s.db.Model(&config).Update("enabled", enabled).Error
}

func (s *SandboxService) GetAllPaths() []models.SandboxPathInfo {
	var paths []SandboxPath
	s.db.Find(&paths)
	result := make([]models.SandboxPathInfo, len(paths))
	for i, p := range paths {
		result[i] = models.SandboxPathInfo{
			ID:          p.ID,
			Path:        p.Path,
			Description: p.Description,
			Enabled:     p.Enabled,
		}
	}
	return result
}
