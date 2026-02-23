package config

import (
	"log"

	"fnchatbot/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InitSystemProviders(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Provider{}).Where("is_system = ?", true).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Printf("System providers already initialized (%d providers found)", count)
		return nil
	}

	log.Printf("Initializing %d system providers...", len(SystemProviders))

	for _, providerDef := range SystemProviders {
		provider := models.Provider{
			ProviderID: providerDef.ProviderID,
			Name:       providerDef.Name,
			Type:       providerDef.Type,
			BaseURL:    providerDef.BaseURL,
			Enabled:    false,
			IsSystem:   true,
			ApiOptions: providerDef.ApiOptions,
		}

		if err := db.Create(&provider).Error; err != nil {
			log.Printf("Failed to create provider %s: %v", providerDef.ProviderID, err)
			continue
		}

		modelDefs := GetModelsByProviderID(providerDef.ProviderID)
		for _, modelDef := range modelDefs {
			model := models.Model{
				ProviderID:             provider.ID,
				ModelID:                modelDef.ModelID,
				Name:                   modelDef.Name,
				Group:                  modelDef.Group,
				Description:            modelDef.Description,
				OwnedBy:                modelDef.OwnedBy,
				Capabilities:           modelDef.Capabilities,
				SupportedEndpointTypes: modelDef.SupportedEndpointTypes,
				EndpointType:           modelDef.EndpointType,
				MaxTokens:              modelDef.MaxTokens,
				InputPrice:             modelDef.InputPrice,
				OutputPrice:            modelDef.OutputPrice,
				Enabled:                true,
			}

			if model.MaxTokens == 0 {
				model.MaxTokens = 4096
			}

			if err := db.Create(&model).Error; err != nil {
				log.Printf("Failed to create model %s for provider %s: %v", modelDef.ModelID, providerDef.ProviderID, err)
			}
		}
	}

	log.Printf("System providers initialized successfully")
	return nil
}

func UpsertSystemProviders(db *gorm.DB) error {
	log.Printf("Upserting %d system providers...", len(SystemProviders))

	for _, providerDef := range SystemProviders {
		var existingProvider models.Provider
		result := db.Where("provider_id = ?", providerDef.ProviderID).First(&existingProvider)

		if result.Error == gorm.ErrRecordNotFound {
			provider := models.Provider{
				ProviderID: providerDef.ProviderID,
				Name:       providerDef.Name,
				Type:       providerDef.Type,
				BaseURL:    providerDef.BaseURL,
				Enabled:    false,
				IsSystem:   true,
				ApiOptions: providerDef.ApiOptions,
			}

			if err := db.Create(&provider).Error; err != nil {
				log.Printf("Failed to create provider %s: %v", providerDef.ProviderID, err)
				continue
			}

			existingProvider = provider
		} else if result.Error != nil {
			log.Printf("Failed to query provider %s: %v", providerDef.ProviderID, result.Error)
			continue
		} else {
			updates := map[string]interface{}{
				"name":        providerDef.Name,
				"type":        providerDef.Type,
				"base_url":    providerDef.BaseURL,
				"api_options": providerDef.ApiOptions,
			}
			if err := db.Model(&existingProvider).Updates(updates).Error; err != nil {
				log.Printf("Failed to update provider %s: %v", providerDef.ProviderID, err)
			}
		}

		if err := syncProviderModels(db, existingProvider.ID, providerDef.ProviderID); err != nil {
			log.Printf("Failed to sync models for provider %s: %v", providerDef.ProviderID, err)
		}
	}

	log.Printf("System providers upserted successfully")
	return nil
}

func syncProviderModels(db *gorm.DB, providerDBID uint, providerID string) error {
	modelDefs := GetModelsByProviderID(providerID)

	var existingModels []models.Model
	if err := db.Where("provider_id = ?", providerDBID).Find(&existingModels).Error; err != nil {
		return err
	}

	existingModelMap := make(map[string]models.Model)
	for _, m := range existingModels {
		existingModelMap[m.ModelID] = m
	}

	for _, modelDef := range modelDefs {
		if existingModel, exists := existingModelMap[modelDef.ModelID]; exists {
			updates := map[string]interface{}{
				"name":                     modelDef.Name,
				"group":                    modelDef.Group,
				"description":              modelDef.Description,
				"owned_by":                 modelDef.OwnedBy,
				"capabilities":             modelDef.Capabilities,
				"supported_endpoint_types": modelDef.SupportedEndpointTypes,
				"endpoint_type":            modelDef.EndpointType,
				"max_tokens":               modelDef.MaxTokens,
				"input_price":              modelDef.InputPrice,
				"output_price":             modelDef.OutputPrice,
			}
			if updates["max_tokens"].(int) == 0 {
				updates["max_tokens"] = 4096
			}
			if err := db.Model(&existingModel).Updates(updates).Error; err != nil {
				log.Printf("Failed to update model %s: %v", modelDef.ModelID, err)
			}
			delete(existingModelMap, modelDef.ModelID)
		} else {
			maxTokens := modelDef.MaxTokens
			if maxTokens == 0 {
				maxTokens = 4096
			}
			model := models.Model{
				ProviderID:             providerDBID,
				ModelID:                modelDef.ModelID,
				Name:                   modelDef.Name,
				Group:                  modelDef.Group,
				Description:            modelDef.Description,
				OwnedBy:                modelDef.OwnedBy,
				Capabilities:           modelDef.Capabilities,
				SupportedEndpointTypes: modelDef.SupportedEndpointTypes,
				EndpointType:           modelDef.EndpointType,
				MaxTokens:              maxTokens,
				InputPrice:             modelDef.InputPrice,
				OutputPrice:            modelDef.OutputPrice,
				Enabled:                true,
			}

			if err := db.Create(&model).Error; err != nil {
				log.Printf("Failed to create model %s: %v", modelDef.ModelID, err)
			}
		}
	}

	return nil
}

func ResetSystemProviders(db *gorm.DB) error {
	if err := db.Where("is_system = ?", true).Delete(&models.Model{}).Error; err != nil {
		return err
	}

	if err := db.Where("is_system = ?", true).Delete(&models.Provider{}).Error; err != nil {
		return err
	}

	return InitSystemProviders(db)
}

func GetProviderWithModels(db *gorm.DB, providerID string) (*models.Provider, error) {
	var provider models.Provider
	if err := db.Where("provider_id = ?", providerID).First(&provider).Error; err != nil {
		return nil, err
	}

	var modelList []models.Model
	if err := db.Where("provider_id = ?", provider.ID).Find(&modelList).Error; err != nil {
		return nil, err
	}

	provider.Models = modelList
	return &provider, nil
}

func GetAllProvidersWithModels(db *gorm.DB) ([]models.Provider, error) {
	var providers []models.Provider
	if err := db.Find(&providers).Error; err != nil {
		return nil, err
	}

	for i := range providers {
		var modelList []models.Model
		if err := db.Where("provider_id = ?", providers[i].ID).Find(&modelList).Error; err != nil {
			return nil, err
		}
		providers[i].Models = modelList
	}

	return providers, nil
}

func UpsertProviderOnConflict(db *gorm.DB, provider *models.Provider) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "type", "base_url", "api_key", "enabled", "api_options", "updated_at"}),
	}).Create(provider).Error
}
