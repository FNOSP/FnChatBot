package api

import (
	"net/http"

	"fnchatbot/internal/db"
	"fnchatbot/internal/models"
	"fnchatbot/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSkills returns all skills
func GetSkills(c *gin.Context) {
	var skills []models.Skill
	if err := db.DB.Order("priority desc, name asc").Find(&skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, skills)
}

// UploadSkill handles skill file upload and parsing
func UploadSkill(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	skill, err := services.ParseSkillFile(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse skill file: " + err.Error()})
		return
	}

	// Check for duplicates
	var existing models.Skill
	if err := db.DB.Where("name = ?", skill.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Skill with this name already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking for duplicates"})
		return
	}

	if err := db.DB.Create(skill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save skill: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, skill)
}

// ToggleSkill updates the enabled status of a skill
func ToggleSkill(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Model(&models.Skill{}).Where("id = ?", id).Update("enabled", input.Enabled).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Skill updated", "enabled": input.Enabled})
}

// DeleteSkill deletes a skill
func DeleteSkill(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Skill{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted"})
}
