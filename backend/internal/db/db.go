package db

import (
	"log"
	"path/filepath"

	"fnchatbot/internal/config"
	"fnchatbot/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

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

	// Auto Migrate the schema
	err = DB.AutoMigrate(
		&models.Provider{},
		&models.Model{},
		&models.ModelConfig{},
		&models.Conversation{},
		&models.Message{},
		&models.MCPConfig{},
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

	// 初始化系统预定义供应商
	if err := config.InitSystemProviders(DB); err != nil {
		log.Printf("Failed to initialize system providers: %v", err)
	}

	log.Println("Database initialized and migrated successfully")
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
		config := models.SandboxConfig{
			Enabled: true,
		}
		if err := DB.Create(&config).Error; err != nil {
			log.Printf("Failed to seed sandbox config: %v", err)
		} else {
			log.Println("Seeded initial sandbox config")
		}
	}
}
