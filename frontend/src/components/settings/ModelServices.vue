<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
import { useTheme } from '../../composables/useTheme'
import { predefinedProviders } from '../../data/providers'
import {
  Layout,
  LayoutSider,
  LayoutContent,
  List,
  ListItem,
  Input,
  Button,
  Table,
  TypographyText,
  TypographyTitle,
  Card,
  Space,
  Switch,
  Modal,
  Checkbox,
  CheckboxGroup,
  Toast,
  Spin,
  Tag,
  Empty
} from '@kousum/semi-ui-vue'
import {
  IconPlus,
  IconDelete,
  IconSetting,
  IconGlobe,
  IconRefresh,
  IconSearch,
  IconServer
} from '@kousum/semi-icons-vue'

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

const API_BASE = 'http://localhost:8080/api'

const loading = ref(false)
const saving = ref(false)
const searchQuery = ref('')
const selectedProviderId = ref<string | null>(null)
const providers = ref<Provider[]>([])

const modelListVisible = ref(false)
const modelListLoading = ref(false)
const availableModels = ref<{ id: string; name: string; owned_by?: string }[]>([])
const selectedModels = ref<string[]>([])

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
  // Normalize API response to avoid empty list when shape differs
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
  }
  )
})

const selectedProvider = computed(() => {
  if (!selectedProviderId.value) return null
  return providers.value.find(p => p.provider_id === selectedProviderId.value)
})

const columns = computed(() => [
  { title: t('settings.displayName'), key: 'name', dataIndex: 'name' },
  { title: t('settings.modelId'), key: 'model_id', dataIndex: 'model_id' },
  { title: t('common.actions'), key: 'actions', width: 120, align: 'right' }
])

const fetchProviders = async () => {
  loading.value = true
  try {
    const res = await axios.get(`${API_BASE}/providers`)
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
    
    if (providers.value.length > 0 && !selectedProviderId.value) {
      const firstEnabled = providers.value.find(p => p.enabled)
      selectedProviderId.value = firstEnabled?.provider_id || providers.value[0].provider_id
    }
  } catch (e) {
    console.error('Failed to fetch providers', e)
    Toast.error(t('common.error'))
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
    await axios.patch(`${API_BASE}/providers/${provider.id}/toggle`)
    Toast.success(provider.enabled ? t('common.enabled') : t('common.disabled'))
  } catch (e) {
    provider.enabled = previousState
    console.error('Failed to toggle provider', e)
    Toast.error(t('common.error'))
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
      await axios.put(`${API_BASE}/providers/${p.id}`, payload)
    } else {
      const res = await axios.post(`${API_BASE}/providers`, payload)
      p.id = res.data.id
    }
    
    Toast.success(t('common.success'))
  } catch (e) {
    console.error('Failed to save provider', e)
    Toast.error(t('common.error'))
  } finally {
    saving.value = false
  }
}

const fetchRemoteModels = async () => {
  if (!selectedProvider.value) return
  
  const p = selectedProvider.value
  if (!p.base_url || !p.api_key) {
    Toast.error(t('settings.modelList.noConfig'))
    return
  }
  
  // Ensure provider exists before fetching remote models
  if (!p.id) {
    await saveProvider()
  }
  if (!p.id) return
  
  modelListLoading.value = true
  try {
    const res = await axios.post(`${API_BASE}/providers/${p.id}/fetch-models`, {
      base_url: p.base_url,
      api_key: p.api_key
    })
    availableModels.value = res.data.models || []
    selectedModels.value = []
    
    if (availableModels.value.length > 0) {
      modelListVisible.value = true
    } else {
      Toast.warning(t('settings.modelList.noModelsFound'))
    }
  } catch (e: any) {
    console.error('Failed to fetch remote models', e)
    Toast.error(e.response?.data?.error || t('settings.modelList.fetchError'))
  } finally {
    modelListLoading.value = false
  }
}

const selectAllModels = () => {
  selectedModels.value = availableModels.value.map(m => m.id)
}

const deselectAllModels = () => {
  selectedModels.value = []
}

const addSelectedModels = async () => {
  if (!selectedProvider.value || !selectedProvider.value.id) return
  
  saving.value = true
  try {
    for (const modelId of selectedModels.value) {
      const modelInfo = availableModels.value.find(m => m.id === modelId)
      await axios.post(`${API_BASE}/providers/${selectedProvider.value.id}/models`, {
        model_id: modelId,
        name: modelInfo?.name || modelId,
        owned_by: modelInfo?.owned_by || '',
        enabled: true
      })
    }
    
    await fetchProviders()
    modelListVisible.value = false
    Toast.success(t('settings.modelList.addSuccess', { count: selectedModels.value.length }))
  } catch (e) {
    console.error('Failed to add models', e)
    Toast.error(t('common.error'))
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
      await axios.delete(`${API_BASE}/models/${model.id}`)
      selectedProvider.value.models = selectedProvider.value.models?.filter(m => m.id !== model.id) || []
      Toast.success(t('common.success'))
    } catch (e) {
      console.error('Failed to delete model', e)
      Toast.error(t('common.error'))
    }
  } else {
    selectedProvider.value.models = selectedProvider.value.models?.filter(m => m !== model) || []
  }
}

const saveModel = async (model: Model) => {
  if (!selectedProvider.value || !selectedProvider.value.id) return
  
  try {
    if (model.id) {
      await axios.put(`${API_BASE}/models/${model.id}`, model)
    } else {
      const payload = {
        ...model,
        provider_id: selectedProvider.value.id
      }
      const res = await axios.post(`${API_BASE}/providers/${selectedProvider.value.id}/models`, payload)
      model.id = res.data.id
    }
    Toast.success(t('common.success'))
  } catch (e) {
    console.error('Failed to save model', e)
    Toast.error(t('common.error'))
  }
}

const getProviderIcon = (providerId: string) => {
  if (providerId === 'openai') return IconGlobe
  if (providerId === 'anthropic') return IconServer
  if (providerId === 'ollama') return IconServer
  return IconSetting
}

onMounted(() => {
  fetchProviders()
})

watch(selectedProviderId, () => {
})
</script>

<template>
  <div class="model-services h-full text-foreground" :class="isDark ? 'bg-zinc-900' : 'bg-white'">
    <Layout class="h-full border rounded-lg overflow-hidden" :class="isDark ? 'bg-zinc-900' : 'bg-white'">
      <LayoutSider class="w-72 border-r flex flex-col" :class="isDark ? 'bg-zinc-900/50' : 'bg-gray-50'">
        <div class="p-4 border-b" :class="isDark ? 'border-zinc-700' : 'border-gray-200'">
          <TypographyTitle :heading="5" class="mb-3">{{ t('settings.providers') }}</TypographyTitle>
          <Input
            v-model="searchQuery"
            :placeholder="t('common.search')"
            :prefix-icon="IconSearch"
            size="small"
            class="w-full"
          />
        </div>
        
        <Spin :spinning="loading">
          <List
            :dataSource="filteredProviders"
            class="flex-1 overflow-y-auto provider-list"
          >
            <template #renderItem="{ item }">
              <ListItem
                :class="[
                  'cursor-pointer transition-colors border-l-2',
                  isDark ? 'hover:bg-white/5 border-transparent' : 'hover:bg-black/5 border-transparent',
                  selectedProviderId === item.provider_id ? (isDark ? 'bg-white/10 border-primary' : 'bg-primary/10 border-primary') : ''
                ]"
                @click="selectedProviderId = item.provider_id"
              >
                <div class="flex items-center justify-between w-full px-4 py-3">
                  <div class="flex items-center gap-3 min-w-0">
                    <div 
                      class="w-8 h-8 rounded flex items-center justify-center flex-shrink-0" 
                      :class="isDark ? 'bg-zinc-800' : 'bg-gray-200'"
                    >
                      <component :is="getProviderIcon(item.provider_id)" class="text-gray-500" />
                    </div>
                    <div class="min-w-0 flex-1">
                      <div class="flex items-center gap-2">
                        <span class="font-medium text-sm truncate">{{ item.name }}</span>
                        <Tag v-if="item.enabled" color="green" size="small">{{ t('common.enabled') }}</Tag>
                      </div>
                      <TypographyText type="secondary" size="small" class="truncate block">
                        {{ item.models?.length || 0 }} {{ t('settings.models') }}
                      </TypographyText>
                    </div>
                  </div>
                  
                  <Switch 
                    size="small" 
                    :checked="item.enabled"
                    @change="() => toggleProvider(item)"
                    @click.stop
                  />
                </div>
              </ListItem>
            </template>
          </List>
        </Spin>
      </LayoutSider>

      <LayoutContent :class="isDark ? 'bg-zinc-900' : 'bg-white'">
        <div v-if="selectedProvider" class="h-full flex flex-col">
          <div class="flex items-center justify-between p-6 border-b" :class="isDark ? 'border-zinc-700' : 'border-gray-200'">
            <div>
              <TypographyTitle :heading="4">
                {{ selectedProvider.name }}
              </TypographyTitle>
              <TypographyText type="secondary" size="small">
                {{ t('settings.configDesc', { name: selectedProvider.name }) }}
              </TypographyText>
            </div>
            <Space>
              <Button 
                theme="solid" 
                type="primary"
                :loading="saving"
                @click="saveProvider"
              >
                {{ t('common.save') }}
              </Button>
            </Space>
          </div>

          <div class="flex-1 overflow-y-auto p-6 space-y-6">
            <Card 
              :title="t('settings.connectionSettings')" 
              :bordered="false" 
              shadow="never" 
              class="border rounded-lg" 
              :class="isDark ? 'bg-zinc-900/50 border-zinc-700' : 'bg-gray-50 border-gray-200'"
            >
              <div class="grid gap-6 max-w-2xl">
                <div class="space-y-2">
                  <label class="text-sm font-medium">{{ t('settings.providerType') }}</label>
                  <Input 
                    :model-value="selectedProvider.type"
                    disabled
                    class="w-full"
                  />
                  <TypographyText type="secondary" size="small">
                    {{ t('settings.providerTypeDesc') }}
                  </TypographyText>
                </div>

                <div class="space-y-2">
                  <label class="text-sm font-medium">{{ t('mcp.baseUrl') }}</label>
                  <Input 
                    v-model="selectedProvider.base_url"
                    placeholder="https://api..."
                    showClear
                    class="w-full"
                  />
                  <TypographyText type="secondary" size="small">
                    {{ t('settings.preview', { url: selectedProvider.base_url || '-' }) }}
                  </TypographyText>
                </div>

                <div class="space-y-2">
                  <label class="text-sm font-medium">{{ t('mcp.apiKey') }}</label>
                  <Input 
                    v-model="selectedProvider.api_key"
                    mode="password"
                    :placeholder="t('mcp.apiKeyPlaceholder')"
                    showClear
                    class="w-full"
                  />
                  <TypographyText type="secondary" size="small">
                    {{ t('settings.apiKeyDesc') }}
                  </TypographyText>
                </div>

                <div class="pt-2 flex gap-2">
                  <Button 
                    :icon="IconRefresh" 
                    :loading="modelListLoading"
                    :disabled="saving"
                    @click="fetchRemoteModels"
                  >
                    {{ t('settings.modelList.fetchBtn') }}
                  </Button>
                </div>
              </div>
            </Card>

            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <TypographyTitle :heading="5">{{ t('settings.models') }}</TypographyTitle>
                <Button 
                  :icon="IconPlus" 
                  @click="addModel"
                  :disabled="!selectedProvider.id"
                >
                  {{ t('settings.addModel') }}
                </Button>
              </div>

              <Table 
                :columns="columns" 
                :dataSource="selectedProvider.models || []" 
                :pagination="false"
                size="small"
                class="border rounded-lg overflow-hidden"
                :class="isDark ? 'border-zinc-700' : 'border-gray-200'"
              >
                <template #empty>
                  <Empty 
                    :description="t('settings.noModels')"
                    class="py-8"
                  />
                </template>

                <template #name="{ record }">
                  <Input 
                    v-model="record.name" 
                    :placeholder="t('settings.displayNamePlaceholder')" 
                    size="small"
                    variant="plain"
                    @blur="saveModel(record)"
                  />
                </template>

                <template #model_id="{ record }">
                  <Input 
                    v-model="record.model_id" 
                    :placeholder="t('settings.modelIdPlaceholder')" 
                    size="small"
                    variant="plain"
                    @blur="saveModel(record)"
                  />
                </template>

                <template #actions="{ record }">
                  <Button 
                    type="danger" 
                    theme="borderless" 
                    :icon="IconDelete" 
                    size="small"
                    @click="removeModel(record)"
                  />
                </template>
              </Table>
            </div>
          </div>
        </div>
        
        <div v-else class="h-full flex items-center justify-center">
          <Empty 
            :description="t('settings.selectProvider')"
            class="text-gray-400"
          />
        </div>
      </LayoutContent>

      <Modal
        v-model:visible="modelListVisible"
        :title="t('settings.modelList.modalTitle')"
        :footer="null"
        width="600px"
      >
        <div class="space-y-4">
          <div class="flex gap-2 items-center">
            <Button size="small" @click="selectAllModels">{{ t('settings.modelList.selectAll') }}</Button>
            <Button size="small" @click="deselectAllModels">{{ t('settings.modelList.deselectAll') }}</Button>
            <TypographyText type="secondary" class="ml-auto">
              {{ t('settings.modelList.selected', { count: selectedModels.length }) }}
            </TypographyText>
          </div>

          <div v-if="availableModels.length === 0" class="text-center py-8 text-gray-500">
            {{ t('settings.modelList.noModelsFound') }}
          </div>

          <CheckboxGroup
            v-else
            v-model="selectedModels"
            direction="vertical"
            class="w-full max-h-80 overflow-y-auto"
          >
            <div
              v-for="model in availableModels"
              :key="model.id"
              class="py-2 border-b last:border-b-0"
              :class="isDark ? 'border-zinc-700' : 'border-gray-100'"
            >
              <Checkbox :value="model.id">
                <div class="flex flex-col">
                  <span class="font-medium">{{ model.name || model.id }}</span>
                  <TypographyText v-if="model.owned_by" type="secondary" size="small">
                    {{ model.owned_by }}
                  </TypographyText>
                </div>
              </Checkbox>
            </div>
          </CheckboxGroup>

          <div class="flex justify-end gap-2 pt-4 border-t" :class="isDark ? 'border-zinc-700' : 'border-gray-200'">
            <Button @click="modelListVisible = false">{{ t('common.cancel') }}</Button>
            <Button 
              theme="solid" 
              type="primary" 
              :disabled="selectedModels.length === 0"
              :loading="saving"
              @click="addSelectedModels"
            >
              {{ t('settings.modelList.addSelected') }}
            </Button>
          </div>
        </div>
      </Modal>
    </Layout>
  </div>
</template>

<style scoped>
:deep(.semi-layout),
:deep(.semi-layout-sider),
:deep(.semi-layout-content) {
  background-color: transparent;
}

:deep(.semi-input-wrapper) {
  background-color: var(--semi-color-fill-0);
}

:deep(.semi-list-item) {
  padding: 0;
}

:deep(.semi-layout-sider) {
  display: flex;
  flex-direction: column;
}

.provider-list {
  flex: 1;
  overflow-y: auto;
}

:deep(.semi-list) {
  height: 100%;
}

:deep(.semi-list-content) {
  height: 100%;
  overflow-y: auto;
}
</style>
