package db

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fnchatbot/internal/config"
	"fnchatbot/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// legacyMCPConfig is used only for migrating from DB to mcp.json.
type legacyMCPConfig struct {
	Name    string `gorm:"column:name"`
	BaseURL string `gorm:"column:base_url"`
	ApiKey  string `gorm:"column:api_key"`
	Enabled bool   `gorm:"column:enabled"`
}

func (legacyMCPConfig) TableName() string {
	return "mcp_configs"
}

var DB *gorm.DB

// InitDB initializes the database connection and performs auto-migration
func InitDB(dbPath string) {
	var err error
	if dbPath == "" {
		dbPath = "fnchatbot.db"
	}

	absPath, _ := filepath.Abs(dbPath)
	log.Printf("Initializing database at: %s", absPath)

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate the schema (MCP config moved to mcp.json, no longer in DB)
	err = DB.AutoMigrate(
		&models.User{},
		&models.Provider{},
		&models.Model{},
		&models.ModelConfig{},
		&models.Session{},
		&models.Message{},
		&models.Part{},
		&models.Skill{},
		&models.AgentTask{},
		&models.SandboxConfig{},
		&models.SandboxPath{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Seed initial data if needed
	seedSkills()
	seedSandboxConfig()

	// Migrate legacy MCP config from DB to mcp.json if that file does not exist
	if err := MigrateMCPConfigToFile("mcp.json"); err != nil {
		log.Printf("MCP config migration: %v", err)
	}

	// 初始化系统预定义供应商
	if err := config.InitSystemProviders(DB); err != nil {
		log.Printf("Failed to initialize system providers: %v", err)
	}
	if err := migrateModelConfigProviders(DB); err != nil {
		log.Printf("Failed to migrate model configs to providers: %v", err)
	}

	log.Println("Database initialized and migrated successfully")
}

type legacyModelConfig struct {
	ID         uint   `gorm:"primaryKey"`
	Provider   string `gorm:"type:varchar(50)"`
	BaseURL    string `gorm:"type:varchar(500)"`
	ApiKey     string `gorm:"type:varchar(500)"`
	ProviderID *uint  `gorm:"index"`
}

func (legacyModelConfig) TableName() string {
	return "model_configs"
}

func migrateModelConfigProviders(db *gorm.DB) error {
	var configs []legacyModelConfig
	if err := db.Find(&configs).Error; err != nil {
		return err
	}

	for _, configRow := range configs {
		if configRow.ProviderID != nil && *configRow.ProviderID > 0 {
			continue
		}
		providerID := strings.TrimSpace(configRow.Provider)
		if providerID == "" {
			continue
		}

		var provider models.Provider
		result := db.Where("provider_id = ?", providerID).First(&provider)
		if result.Error != nil {
			if result.Error != gorm.ErrRecordNotFound {
				return result.Error
			}
			def := config.GetProviderByID(providerID)
			name := providerID
			providerType := models.ProviderTypeOpenAI
			baseURL := configRow.BaseURL
			apiOptions := models.ProviderApiOptions{}
			isSystem := false
			if def != nil {
				name = def.Name
				providerType = def.Type
				apiOptions = def.ApiOptions
				isSystem = true
				if baseURL == "" {
					baseURL = def.BaseURL
				}
			}
			provider = models.Provider{
				ProviderID: providerID,
				Name:       name,
				Type:       providerType,
				BaseURL:    baseURL,
				APIKey:     configRow.ApiKey,
				Enabled:    configRow.ApiKey != "",
				IsSystem:   isSystem,
				ApiOptions: apiOptions,
			}
			if err := db.Create(&provider).Error; err != nil {
				return err
			}
		} else {
			updates := map[string]interface{}{}
			if provider.BaseURL == "" && configRow.BaseURL != "" {
				updates["base_url"] = configRow.BaseURL
			}
			if provider.APIKey == "" && configRow.ApiKey != "" {
				updates["api_key"] = configRow.ApiKey
				updates["enabled"] = true
			}
			if len(updates) > 0 {
				if err := db.Model(&provider).Updates(updates).Error; err != nil {
					return err
				}
			}
		}

		if err := db.Model(&models.ModelConfig{}).
			Where("id = ?", configRow.ID).
			Update("provider_id", provider.ID).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedSkills() {
	var count int64
	DB.Model(&models.Skill{}).Count(&count)
	if count == 0 {
		skills := []models.Skill{
			{Name: "web_search", Description: "Search the web for information", Priority: 10},
			{Name: "code_execute", Description: "Execute code snippets", Priority: 5},
			{Name: "file_read", Description: "Read files from the local system", Priority: 8},
			{Name: "calculator", Description: "Perform mathematical calculations", Priority: 3},
		}
		DB.Create(&skills)
		log.Println("Seeded initial skills")
	}
}

func seedSandboxConfig() {
	var count int64
	DB.Model(&models.SandboxConfig{}).Count(&count)
	if count == 0 {
		sandboxConfig := models.SandboxConfig{
			Enabled: true,
		}
		if err := DB.Create(&sandboxConfig).Error; err != nil {
			log.Printf("Failed to seed sandbox config: %v", err)
		} else {
			log.Println("Seeded initial sandbox config")
		}
	}
}

// MigrateMCPConfigToFile exports legacy mcp_configs rows to mcp.json if the file does not exist.
func MigrateMCPConfigToFile(mcpFilePath string) error {
	if mcpFilePath == "" {
		mcpFilePath = "mcp.json"
	}
	if _, err := os.Stat(mcpFilePath); err == nil {
		return nil // file exists, skip
	}
	var rows []legacyMCPConfig
	if err := DB.Find(&rows).Error; err != nil {
		// Table might not exist on fresh install
		return nil
	}
	if len(rows) == 0 {
		return nil
	}
	servers := make(map[string]models.MCPServerConfig)
	for _, r := range rows {
		name := strings.TrimSpace(r.Name)
		if name == "" {
			continue
		}
		servers[name] = models.MCPServerConfig{
			Type:    models.MCPTypeRemote,
			URL:     r.BaseURL,
			ApiKey:  r.ApiKey,
			Enabled: r.Enabled,
			Timeout: 5000,
		}
	}
	if len(servers) == 0 {
		return nil
	}
	data, err := json.MarshalIndent(&models.MCPFile{Servers: servers}, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(mcpFilePath, data, 0644); err != nil {
		return err
	}
	log.Printf("Migrated %d MCP config(s) from DB to %s", len(servers), mcpFilePath)
	return nil
}
