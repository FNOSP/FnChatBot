<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { http } from '../../services/http'
import { useI18n } from 'vue-i18n'
import { useTheme } from '../../composables/useTheme'
import { predefinedProviders } from '../../data/providers'
import { MessagePlugin } from 'tdesign-vue-next'
import {
  AddIcon,
  DeleteIcon,
  SettingIcon,
  InternetIcon,
  RefreshIcon,
  SearchIcon,
  ServerIcon
} from 'tdesign-icons-vue-next'

const { t } = useI18n()
const { isDark } = useTheme()

interface Model {
  id?: number
  model_id: string
  name: string
  provider_id?: number
  group?: string
  description?: string
  owned_by?: string
  capabilities?: string[]
  max_tokens?: number
  input_price?: number
  output_price?: number
  enabled?: boolean
  is_default?: boolean
}

interface Provider {
  id?: number
  provider_id: string
  name: string
  type: string
  base_url: string
  api_key?: string
  enabled: boolean
  is_system?: boolean
  models?: Model[]
}

const loading = ref(false)
const saving = ref(false)
const searchQuery = ref('')
const selectedProviderId = ref<string | null>(null)
const providers = ref<Provider[]>([])

const modelListVisible = ref(false)
const modelListLoading = ref(false)
const availableModels = ref<{ id: string; name: string; owned_by?: string }[]>([])

const buildDefaultProviders = (): Provider[] => {
  return predefinedProviders.map(predefined => ({
    provider_id: predefined.id,
    name: predefined.name,
    type: predefined.type,
    base_url: predefined.baseUrl,
    enabled: false,
    models: []
  }))
}

const normalizeProvidersResponse = (data: any): Provider[] => {
  if (Array.isArray(data)) return data
  if (Array.isArray(data?.providers)) return data.providers
  if (Array.isArray(data?.data)) return data.data
  return []
}

const filteredProviders = computed(() => {
  if (!searchQuery.value) return providers.value
  const query = searchQuery.value.toLowerCase()
  return providers.value.filter(p => {
    const name = (p.name || '').toLowerCase()
    const providerId = (p.provider_id || '').toLowerCase()
    return name.includes(query) || providerId.includes(query)
  })
})

const selectedProvider = computed(() => {
  if (!selectedProviderId.value) return null
  return providers.value.find(p => p.provider_id === selectedProviderId.value)
})

const localModelMap = computed(() => {
  const map = new Map<string, Model>()
  const models = selectedProvider.value?.models || []
  for (const model of models) {
    if (model.model_id) {
      map.set(model.model_id, model)
    }
  }
  return map
})

const sortProvidersByName = (list: Provider[]) => {
  return [...list].sort((a, b) => (a.name || '').localeCompare(b.name || '', undefined, { sensitivity: 'base' }))
}

const getDefaultProviderId = (list: Provider[]) => {
  const enabledProviders = list.filter(p => p.enabled)
  if (enabledProviders.length > 0) {
    return sortProvidersByName(enabledProviders)[0].provider_id
  }
  const sortedAll = sortProvidersByName(list)
  return sortedAll[0]?.provider_id || null
}

const columns = computed(() => [
  { title: t('settings.displayName'), colKey: 'name', width: 200 },
  { title: t('settings.modelId'), colKey: 'model_id', width: 200 },
  { title: t('common.actions'), colKey: 'actions', width: 100, align: 'right' }
])

const fetchProviders = async () => {
  loading.value = true
  try {
    const res = await http.get('/providers')
    const backendProviders = normalizeProvidersResponse(res.data)
    
    const mergedProviders: Provider[] = predefinedProviders.map(predefined => {
      const existing = backendProviders.find(bp => bp.provider_id === predefined.id)
      if (existing) {
        return {
          ...existing,
          name: predefined.name,
          type: predefined.type,
          base_url: existing.base_url || predefined.baseUrl
        }
      }
      return {
        provider_id: predefined.id,
        name: predefined.name,
        type: predefined.type,
        base_url: predefined.baseUrl,
        enabled: false,
        models: []
      }
    })
    
    const customProviders = backendProviders.filter(
      bp => !predefinedProviders.find(pp => pp.id === bp.provider_id)
    )
    
    const mergedList = [...mergedProviders, ...customProviders]
    providers.value = mergedList.length > 0 ? mergedList : buildDefaultProviders()
    
    if (providers.value.length > 0) {
      const fallbackId = getDefaultProviderId(providers.value)
      const hasSelected = selectedProviderId.value && providers.value.some(p => p.provider_id === selectedProviderId.value)
      if (!hasSelected && fallbackId) {
        selectedProviderId.value = fallbackId
      }
    }
  } catch (e) {
    console.error('Failed to fetch providers', e)
    MessagePlugin.error(t('common.error'))
  } finally {
    loading.value = false
  }
}

const toggleProvider = async (provider: Provider) => {
  if (!provider.id) {
    provider.enabled = !provider.enabled
    return
  }
  
  const previousState = provider.enabled
  provider.enabled = !provider.enabled
  
  try {
    await http.patch(`/providers/${provider.id}/toggle`)
    MessagePlugin.success(provider.enabled ? t('common.enabled') : t('common.disabled'))
  } catch (e) {
    provider.enabled = previousState
    console.error('Failed to toggle provider', e)
    MessagePlugin.error(t('common.error'))
  }
}

const saveProvider = async () => {
  if (!selectedProvider.value) return
  if (saving.value) return
  
  saving.value = true
  const p = selectedProvider.value
  
  try {
    const payload = {
      provider_id: p.provider_id,
      name: p.name,
      type: p.type,
      base_url: p.base_url,
      api_key: p.api_key || '',
      enabled: p.enabled
    }
    
    if (p.id) {
      await http.put(`/providers/${p.id}`, payload)
    } else {
      const res = await http.post('/providers', payload)
      p.id = res.data.id
    }
    
    MessagePlugin.success(t('common.success'))
  } catch (e) {
    console.error('Failed to save provider', e)
    MessagePlugin.error(t('common.error'))
  } finally {
    saving.value = false
  }
}

const fetchRemoteModels = async () => {
  if (!selectedProvider.value) return
  
  const p = selectedProvider.value
  if (!p.base_url || !p.api_key) {
    MessagePlugin.error(t('settings.modelList.noConfig'))
    return
  }
  
  if (!p.id) {
    await saveProvider()
  }
  if (!p.id) return
  
  modelListLoading.value = true
  try {
    const res = await http.post(`/providers/${p.id}/fetch-models`, {
      base_url: p.base_url,
      api_key: p.api_key
    })
    availableModels.value = res.data.models || []
    
    if (availableModels.value.length > 0) {
      modelListVisible.value = true
    } else {
      MessagePlugin.warning(t('settings.modelList.noModelsFound'))
    }
  } catch (e: any) {
    console.error('Failed to fetch remote models', e)
    MessagePlugin.error(e.response?.data?.error || t('settings.modelList.fetchError'))
  } finally {
    modelListLoading.value = false
  }
}

const isModelAdded = (modelId: string) => {
  return localModelMap.value.has(modelId)
}

const addRemoteModel = async (model: { id: string; name: string; owned_by?: string }) => {
  if (!selectedProvider.value || !selectedProvider.value.id) return
  if (isModelAdded(model.id)) return

  saving.value = true
  try {
    const res = await http.post(`/providers/${selectedProvider.value.id}/models`, {
      model_id: model.id,
      name: model.name || model.id,
      owned_by: model.owned_by || '',
      enabled: true
    })
    if (!selectedProvider.value.models) {
      selectedProvider.value.models = []
    }
    selectedProvider.value.models.push(res.data)
    MessagePlugin.success(t('common.success'))
  } catch (e) {
    console.error('Failed to add model', e)
    MessagePlugin.error(t('common.error'))
  } finally {
    saving.value = false
  }
}

const removeRemoteModel = async (modelId: string) => {
  if (!selectedProvider.value) return
  const target = localModelMap.value.get(modelId)
  if (!target?.id) return

  saving.value = true
  try {
    await http.delete(`/models/${target.id}`)
    selectedProvider.value.models = selectedProvider.value.models?.filter(m => m.id !== target.id) || []
    MessagePlugin.success(t('common.success'))
  } catch (e) {
    console.error('Failed to delete model', e)
    MessagePlugin.error(t('common.error'))
  } finally {
    saving.value = false
  }
}

const addModel = () => {
  if (!selectedProvider.value) return
  if (!selectedProvider.value.models) {
    selectedProvider.value.models = []
  }
  selectedProvider.value.models.push({
    model_id: '',
    name: '',
    enabled: true
  })
}

const removeModel = async (model: Model) => {
  if (!selectedProvider.value) return
  
  if (model.id) {
    try {
      await http.delete(`/models/${model.id}`)
      selectedProvider.value.models = selectedProvider.value.models?.filter(m => m.id !== model.id) || []
      MessagePlugin.success(t('common.success'))
    } catch (e) {
      console.error('Failed to delete model', e)
      MessagePlugin.error(t('common.error'))
    }
  } else {
    selectedProvider.value.models = selectedProvider.value.models?.filter(m => m !== model) || []
  }
}

const saveModel = async (model: Model) => {
  if (!selectedProvider.value || !selectedProvider.value.id) return
  
  try {
    if (model.id) {
      await http.put(`/models/${model.id}`, model)
    } else {
      const payload = {
        ...model,
        provider_id: selectedProvider.value.id
      }
      const res = await http.post(`/providers/${selectedProvider.value.id}/models`, payload)
      model.id = res.data.id
    }
    MessagePlugin.success(t('common.success'))
  } catch (e) {
    console.error('Failed to save model', e)
    MessagePlugin.error(t('common.error'))
  }
}

const getProviderIcon = (providerId: string) => {
  if (providerId === 'openai') return InternetIcon
  if (providerId === 'anthropic') return ServerIcon
  if (providerId === 'ollama') return ServerIcon
  return SettingIcon
}

onMounted(() => {
  fetchProviders()
})
</script>

<template>
  <div class="h-full text-text-primary bg-bg-primary">
    <t-layout class="h-full border border-border rounded-lg overflow-hidden bg-bg-card">
      <t-aside
        class="border-r border-border flex flex-col bg-bg-secondary"
        width="18rem"
      >
        <div class="p-4 border-b border-border">
          <h3 class="text-lg font-bold mb-3">{{ t('settings.providers') }}</h3>
          <t-input
            v-model="searchQuery"
            :placeholder="t('common.search')"
            size="small"
            class="w-full"
          >
            <template #prefix-icon>
              <SearchIcon />
            </template>
          </t-input>
        </div>
        
        <div class="flex-1 min-h-0 overflow-y-auto">
          <t-loading :loading="loading" class="w-full h-full">
            <t-list class="provider-list" :split="true">
              <t-list-item
                v-for="item in filteredProviders"
                :key="item.provider_id"
                :class="[
                  'cursor-pointer transition-colors border-l-2',
                  'hover:bg-bg-hover border-transparent',
                  selectedProviderId === item.provider_id ? 'bg-brand-light border-brand' : ''
                ]"
                @click="selectedProviderId = item.provider_id"
              >
                <div class="flex items-center justify-between w-full px-4 py-3">
                  <div class="flex items-center gap-3 min-w-0">
                    <div 
                      class="w-8 h-8 rounded flex items-center justify-center flex-shrink-0 bg-bg-secondary-active" 
                    >
                      <component :is="getProviderIcon(item.provider_id)" class="text-text-secondary" />
                    </div>
                    <div class="min-w-0 flex-1">
                      <div class="flex items-center gap-2">
                        <span class="font-medium text-sm truncate text-text-primary">{{ item.name }}</span>
                        <t-tag v-if="item.enabled" theme="success" size="small" variant="light">ON</t-tag>
                      </div>
                      <span class="text-xs text-text-secondary truncate block">
                        {{ item.models?.length || 0 }} {{ t('settings.models') }}
                      </span>
                    </div>
                  </div>
                  
                  <t-switch 
                    size="small" 
                    :value="item.enabled"
                    @change="() => toggleProvider(item)"
                    @click.stop
                  />
                </div>
              </t-list-item>
            </t-list>
          </t-loading>
        </div>
      </t-aside>

      <t-content class="bg-bg-primary">
        <div v-if="selectedProvider" class="h-full flex flex-col">
          <div class="flex items-center justify-between p-6 border-b border-border">
            <div>
              <h4 class="text-lg font-bold">
                {{ selectedProvider.name }}
              </h4>
              <span class="text-sm text-text-secondary">
                {{ t('settings.configDesc', { name: selectedProvider.name }) }}
              </span>
            </div>
            <div class="flex gap-2">
              <t-button 
                theme="primary"
                :loading="saving"
                @click="saveProvider"
              >
                {{ t('common.save') }}
              </t-button>
            </div>
          </div>

          <div class="flex-1 overflow-y-auto p-6 space-y-6">
            <t-card 
              :title="t('settings.connectionSettings')" 
              :bordered="true" 
              class="rounded-lg bg-bg-card" 
            >
              <div class="grid gap-6 max-w-2xl">
                <div class="space-y-2">
                  <label class="text-sm font-medium">{{ t('settings.providerType') }}</label>
                  <t-input 
                    :value="selectedProvider.type"
                    disabled
                    class="w-full"
                  />
                  <span class="text-xs text-text-secondary">
                    {{ t('settings.providerTypeDesc') }}
                  </span>
                </div>

                <div class="space-y-2">
                  <label class="text-sm font-medium">{{ t('mcp.baseUrl') }}</label>
                  <t-input 
                    v-model="selectedProvider.base_url"
                    placeholder="https://api..."
                    clearable
                    class="w-full"
                  />
                  <span class="text-xs text-text-secondary">
                    {{ t('settings.preview', { url: selectedProvider.base_url || '-' }) }}
                  </span>
                </div>

                <div class="space-y-2">
                  <label class="text-sm font-medium">{{ t('mcp.apiKey') }}</label>
                  <t-input 
                    v-model="selectedProvider.api_key"
                    type="password"
                    :placeholder="t('mcp.apiKeyPlaceholder')"
                    clearable
                    class="w-full"
                  />
                  <span class="text-xs text-text-secondary">
                    {{ t('settings.apiKeyDesc') }}
                  </span>
                </div>

                <div class="pt-2 flex gap-2">
                  <t-button 
                    variant="outline"
                    :loading="modelListLoading"
                    :disabled="saving"
                    @click="fetchRemoteModels"
                  >
                    <template #icon><RefreshIcon /></template>
                    {{ t('settings.modelList.fetchBtn') }}
                  </t-button>
                </div>
              </div>
            </t-card>

            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h5 class="text-base font-bold">{{ t('settings.models') }}</h5>
                <t-button 
                  variant="text"
                  @click="addModel"
                  :disabled="!selectedProvider.id"
                >
                  <template #icon><AddIcon /></template>
                  {{ t('settings.addModel') }}
                </t-button>
              </div>

              <t-table 
                :columns="columns" 
                :data="selectedProvider.models || []" 
                :pagination="null"
                size="small"
                row-key="id"
                class="border border-border rounded-lg overflow-hidden"
              >
                <template #empty>
                  <div class="p-8 text-center text-text-muted">
                    {{ t('settings.noModels') }}
                  </div>
                </template>

                <template #name="{ row }">
                  <t-input 
                    v-model="row.name" 
                    :placeholder="t('settings.displayNamePlaceholder')" 
                    size="small"
                    borderless
                    @blur="saveModel(row)"
                  />
                </template>

                <template #model_id="{ row }">
                  <t-input 
                    v-model="row.model_id" 
                    :placeholder="t('settings.modelIdPlaceholder')" 
                    size="small"
                    borderless
                    @blur="saveModel(row)"
                  />
                </template>

                <template #actions="{ row }">
                  <t-button 
                    theme="danger" 
                    variant="text" 
                    shape="square"
                    size="small"
                    @click="removeModel(row)"
                  >
                    <template #icon><DeleteIcon /></template>
                  </t-button>
                </template>
              </t-table>
            </div>
          </div>
        </div>
        
        <div v-else class="h-full flex items-center justify-center text-text-muted">
          {{ t('settings.selectProvider') }}
        </div>
      </t-content>

      <t-dialog
        v-model:visible="modelListVisible"
        :header="t('settings.modelList.modalTitle')"
        :footer="false"
        width="600px"
      >
        <div class="space-y-4">
          <div v-if="availableModels.length === 0" class="text-center py-8 text-text-muted">
            {{ t('settings.modelList.noModelsFound') }}
          </div>

          <div
            v-for="model in availableModels"
            :key="model.id"
            class="py-3 border-b border-border last:border-b-0"
          >
            <div class="flex items-center justify-between gap-4">
              <div class="flex flex-col min-w-0">
                <span class="font-medium truncate text-text-primary">{{ model.name || model.id }}</span>
                <span v-if="model.owned_by" class="text-xs text-text-secondary truncate">
                  {{ model.owned_by }}
                </span>
              </div>
              <div class="flex items-center gap-2">
                <t-tag v-if="isModelAdded(model.id)" theme="success" size="small">ON</t-tag>
                <t-button
                  v-if="isModelAdded(model.id)"
                  theme="danger"
                  variant="text"
                  shape="square"
                  size="small"
                  :disabled="saving"
                  @click="removeRemoteModel(model.id)"
                >
                  <template #icon><DeleteIcon /></template>
                </t-button>
                <t-button
                  v-else
                  variant="text"
                  shape="square"
                  size="small"
                  :disabled="saving"
                  @click="addRemoteModel(model)"
                >
                  <template #icon><AddIcon /></template>
                </t-button>
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-4 border-t border-border">
            <t-button variant="outline" @click="modelListVisible = false">{{ t('common.cancel') }}</t-button>
          </div>
        </div>
      </t-dialog>
    </t-layout>
  </div>
</template>

<style scoped>
.provider-list {
  flex: 1;
  overflow-y: auto;
}
</style>
