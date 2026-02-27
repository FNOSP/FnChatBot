<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { http } from '../../services/http'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'
import { 
  AddIcon, 
  EditIcon, 
  DeleteIcon, 
  ServerIcon,
  PlayCircleIcon,
  StopCircleIcon
} from 'tdesign-icons-vue-next'

const { t } = useI18n()

// --- Types ---
interface MCPConfig {
  id?: number
  name: string
  base_url: string
  api_key: string
  enabled: boolean
}

// --- State ---
const servers = ref<MCPConfig[]>([])
const loading = ref(false)
const showDialog = ref(false)
const editingServer = ref<MCPConfig | null>(null)
const form = ref<MCPConfig>({
  name: '',
  base_url: '',
  api_key: '',
  enabled: true
})
const saving = ref(false)

// --- Actions ---
// Load all configured MCP servers from backend
const fetchServers = async () => {
  loading.value = true
  try {
    const res = await http.get('/mcp')
    servers.value = res.data
  } catch (e) {
    console.error('Failed to fetch MCP servers', e)
    MessagePlugin.error(t('common.error'))
  } finally {
    loading.value = false
  }
}

// Open dialog for creating a new MCP server
const openAddDialog = () => {
  editingServer.value = null
  form.value = { name: '', base_url: '', api_key: '', enabled: true }
  showDialog.value = true
}

// Open dialog for editing an existing MCP server
const openEditDialog = (server: MCPConfig) => {
  editingServer.value = server
  form.value = { ...server }
  showDialog.value = true
}

// Delete a single MCP server by id
const deleteServer = async (id: number) => {
  const confirm = window.confirm(t('common.deleteConfirm'))
  if (!confirm) return
  try {
    await http.delete(`/mcp/${id}`)
    await fetchServers()
    MessagePlugin.success(t('common.success'))
  } catch (e) {
    console.error('Failed to delete server', e)
    MessagePlugin.error(t('common.error'))
  }
}

// Persist MCP server changes (create or update)
const saveServer = async () => {
  saving.value = true
  try {
    if (editingServer.value && editingServer.value.id) {
      await http.put(`/mcp/${editingServer.value.id}`, form.value)
    } else {
      await http.post('/mcp', form.value)
    }
    await fetchServers()
    showDialog.value = false
    MessagePlugin.success(t('common.success'))
  } catch (e) {
    console.error('Failed to save server', e)
    MessagePlugin.error(t('common.error'))
  } finally {
    saving.value = false
  }
}

// Toggle MCP server enabled flag with optimistic UI update
const toggleServer = async (server: MCPConfig) => {
    if (!server.id) return
    // Optimistic update
    const newValue = !server.enabled
    server.enabled = newValue
    
    try {
        await http.put(`/mcp/${server.id}`, { ...server, enabled: newValue })
        MessagePlugin.success(server.enabled ? t('common.enabled') : t('common.disabled'))
    } catch (e) {
        // Revert on error
        server.enabled = !newValue
        console.error('Failed to toggle server', e)
        MessagePlugin.error(t('common.error'))
    }
}

onMounted(() => {
  fetchServers()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="flex items-center gap-2 text-lg font-bold">
            <ServerIcon />
            {{ t('mcp.title') }}
        </h3>
        <span class="text-text-secondary">
            {{ t('mcp.description') }}
        </span>
      </div>
      <t-button 
        theme="primary" 
        @click="openAddDialog"
      >
        <template #icon><AddIcon /></template>
        {{ t('mcp.addServer') }}
      </t-button>
    </div>

    <!-- List -->
    <t-loading :loading="loading">
        <div v-if="servers.length === 0 && !loading" class="text-center py-12">
            <div class="flex flex-col items-center gap-4">
                <ServerIcon size="48" class="text-text-muted" />
                <h4 class="text-lg font-bold">{{ t('mcp.noServers') }}</h4>
                <p class="text-text-secondary">{{ t('mcp.noServersDesc') }}</p>
                <t-button @click="openAddDialog">{{ t('mcp.addFirst') }}</t-button>
            </div>
        </div>

        <t-list
            v-else
            :split="true"
            class="bg-bg-card rounded-lg border border-border"
        >
            <t-list-item v-for="item in servers" :key="item.id" class="hover:bg-bg-hover transition-colors">
                <div class="flex items-center justify-between w-full">
                    <div class="flex items-center gap-3">
                        <div class="p-2 bg-bg-secondary rounded">
                            <ServerIcon />
                        </div>
                        <div class="flex flex-col">
                            <div class="flex items-center gap-2">
                                <span class="font-medium text-text-primary">{{ item.name }}</span>
                                <t-tag :theme="item.enabled ? 'success' : 'default'" size="small">
                                    {{ item.enabled ? t('common.enabled') : t('common.disabled') }}
                                </t-tag>
                            </div>
                            <span class="text-xs text-text-secondary">{{ item.base_url }}</span>
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <t-button
                            variant="text"
                            shape="square"
                            :theme="item.enabled ? 'success' : 'default'"
                            @click="toggleServer(item)"
                        >
                            <template #icon>
                                <PlayCircleIcon v-if="item.enabled" />
                                <StopCircleIcon v-else />
                            </template>
                        </t-button>
                        <t-button
                            variant="text"
                            shape="square"
                            @click="openEditDialog(item)"
                        >
                            <template #icon><EditIcon /></template>
                        </t-button>
                        <t-button
                            variant="text"
                            shape="square"
                            theme="danger"
                            @click="deleteServer(item.id!)"
                        >
                            <template #icon><DeleteIcon /></template>
                        </t-button>
                    </div>
                </div>
            </t-list-item>
        </t-list>
    </t-loading>

    <!-- Dialog -->
    <t-dialog
        v-model:visible="showDialog"
        :header="editingServer ? t('mcp.editServer') : t('mcp.addServer')"
        :confirm-btn="t('common.save')"
        :cancel-btn="t('common.cancel')"
        @confirm="saveServer"
        @cancel="showDialog = false"
        :confirm-loading="saving"
        :confirm-btn-props="{ disabled: !form.name || !form.base_url }"
    >
        <div class="space-y-4">
            <div class="grid gap-2">
                <label class="text-sm font-medium">{{ t('mcp.name') }}</label>
                <t-input 
                    v-model="form.name"
                    :placeholder="t('mcp.namePlaceholder')"
                    autofocus
                />
            </div>
            <div class="grid gap-2">
                <label class="text-sm font-medium">{{ t('mcp.baseUrl') }}</label>
                <t-input 
                    v-model="form.base_url"
                    :placeholder="t('mcp.urlPlaceholder')"
                />
            </div>
            <div class="grid gap-2">
                <label class="text-sm font-medium">{{ t('mcp.apiKey') }}</label>
                <t-input 
                    v-model="form.api_key"
                    type="password"
                    :placeholder="t('mcp.apiKeyPlaceholder')"
                />
            </div>
        </div>
    </t-dialog>
  </div>
</template>
