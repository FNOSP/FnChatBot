package adapters

import (
	"sync"

	"fnchatbot/internal/models"
	"fnchatbot/internal/services"
)

var (
	adapters     map[models.ProviderType]services.ProviderAdapter
	adaptersOnce sync.Once
)

func RegisterAdapter(providerType models.ProviderType, adapter services.ProviderAdapter) {
	if adapters == nil {
		adapters = make(map[models.ProviderType]services.ProviderAdapter)
	}
	adapters[providerType] = adapter
}

func GetAdapter(providerType models.ProviderType) services.ProviderAdapter {
	adaptersOnce.Do(initAdapters)
	return adapters[providerType]
}

func initAdapters() {
	adapters = make(map[models.ProviderType]services.ProviderAdapter)

	adapters[models.ProviderTypeOpenAI] = NewOpenAIAdapter()
	adapters[models.ProviderTypeOpenAIResponse] = NewOpenAIAdapter()
	adapters[models.ProviderTypeAnthropic] = NewAnthropicAdapter()
	adapters[models.ProviderTypeOllama] = NewOllamaAdapter()
	adapters[models.ProviderTypeGemini] = NewGeminiAdapter()
	adapters[models.ProviderTypeAzureOpenAI] = NewAzureAdapter()

	adapters[models.ProviderTypeMistral] = NewOpenAIAdapter()
	adapters[models.ProviderTypeNewAPI] = NewOpenAIAdapter()
	adapters[models.ProviderTypeGateway] = NewOpenAIAdapter()
	adapters[models.ProviderTypeVertexAI] = NewOpenAIAdapter()
	adapters[models.ProviderTypeAWBedrock] = NewOpenAIAdapter()
	adapters[models.ProviderTypeVertexAnthropic] = NewAnthropicAdapter()
}

func GetAdapterForProvider(provider *models.Provider) services.ProviderAdapter {
	if provider == nil {
		return nil
	}
	return GetAdapter(provider.Type)
}
