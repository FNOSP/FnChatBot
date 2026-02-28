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
  StopCircleIcon,
  LoadingIcon,
} from 'tdesign-icons-vue-next'

const { t } = useI18n()

// --- Types ---
type MCPType = 'local' | 'remote'
type MCPStatusType = 'connected' | 'disabled' | 'failed' | 'unknown'

interface MCPServerInfo {
  name: string
  type: MCPType
  command?: string[]
  env?: Record<string, string>
  url?: string
  api_key?: string
  headers?: Record<string, string>
  enabled: boolean
  timeout?: number
  status: MCPStatusType
  error?: string
}

// Form uses command as string (space-separated) for editing; backend expects command[].
const defaultForm = (): Partial<MCPServerInfo> & { name: string; type: MCPType; enabled: boolean; commandStr?: string } => ({
  name: '',
  type: 'remote',
  url: '',
  api_key: '',
  enabled: true,
  status: 'unknown',
  commandStr: '',
})

// --- State ---
const servers = ref<MCPServerInfo[]>([])
const loading = ref(false)
const showDialog = ref(false)
const editingServer = ref<MCPServerInfo | null>(null)
const form = ref(defaultForm())
const saving = ref(false)
const checkingName = ref<string | null>(null)

// --- Actions ---
const fetchServers = async () => {
  loading.value = true
  try {
    const res = await http.get('/mcp')
    servers.value = res.data ?? []
  } catch (e) {
    console.error('Failed to fetch MCP servers', e)
    MessagePlugin.error(t('common.error'))
  } finally {
    loading.value = false
  }
}

const openAddDialog = () => {
  editingServer.value = null
  form.value = { ...defaultForm(), name: '', type: 'remote', enabled: true, commandStr: '' }
  showDialog.value = true
}

const openEditDialog = (server: MCPServerInfo) => {
  editingServer.value = server
  const commandStr = Array.isArray(server.command) ? server.command.join(' ') : ''
  form.value = {
    name: server.name,
    type: server.type,
    url: server.url ?? '',
    api_key: server.api_key ?? '',
    enabled: server.enabled,
    status: server.status,
    commandStr,
  }
  showDialog.value = true
}

const deleteServer = async (server: MCPServerInfo) => {
  if (!window.confirm(t('common.deleteConfirm'))) return
  try {
    await http.delete(`/mcp/${encodeURIComponent(server.name)}`)
    await fetchServers()
    MessagePlugin.success(t('common.success'))
  } catch (e) {
    console.error('Failed to delete server', e)
    MessagePlugin.error(t('common.error'))
  }
}

const buildPayload = () => {
  const type = form.value.type!
  const payload: Record<string, unknown> = {
    name: form.value.name,
    type,
    enabled: form.value.enabled,
  }
  if (form.value.timeout) payload.timeout = form.value.timeout
  if (type === 'local') {
    const cmd = (form.value.commandStr ?? '').trim().split(/\s+/).filter(Boolean)
    payload.command = cmd
    if (form.value.env && Object.keys(form.value.env).length) payload.env = form.value.env
  } else {
    payload.url = form.value.url ?? ''
    payload.api_key = form.value.api_key ?? ''
    if (form.value.headers && Object.keys(form.value.headers).length) payload.headers = form.value.headers
  }
  return payload
}

const saveServer = async () => {
  if (!form.value.name?.trim()) return
  saving.value = true
  try {
    const payload = buildPayload()
    if (editingServer.value) {
      await http.put(`/mcp/${encodeURIComponent(editingServer.value.name)}`, payload)
    } else {
      await http.post('/mcp', payload)
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

const toggleServer = async (server: MCPServerInfo) => {
  const newValue = !server.enabled
  server.enabled = newValue
  try {
    const payload: Record<string, unknown> = {
      type: server.type,
      enabled: newValue,
      timeout: server.timeout,
    }
    if (server.type === 'local') {
      payload.command = server.command ?? []
      payload.env = server.env
    } else {
      payload.url = server.url ?? ''
      payload.api_key = server.api_key ?? ''
      payload.headers = server.headers
    }
    await http.put(`/mcp/${encodeURIComponent(server.name)}`, payload)
    MessagePlugin.success(server.enabled ? t('common.enabled') : t('common.disabled'))
    if (newValue) {
      checkingName.value = server.name
      try {
        await http.post(`/mcp/${encodeURIComponent(server.name)}/check`)
      } finally {
        checkingName.value = null
      }
      await fetchServers()
    }
  } catch (e) {
    server.enabled = !newValue
    console.error('Failed to toggle server', e)
    MessagePlugin.error(t('common.error'))
  }
}

const statusTheme = (status: MCPStatusType) => {
  switch (status) {
    case 'connected': return 'success'
    case 'failed': return 'danger'
    case 'disabled': return 'default'
    default: return 'warning'
  }
}

const statusLabel = (status: MCPStatusType) => {
  switch (status) {
    case 'connected': return t('mcp.statusConnected')
    case 'failed': return t('mcp.statusFailed')
    case 'disabled': return t('mcp.statusDisabled')
    default: return t('mcp.statusUnknown')
  }
}

const displaySubtitle = (item: MCPServerInfo) => {
  if (item.type === 'local' && Array.isArray(item.command)) return item.command.join(' ')
  return item.url ?? ''
}

const canSave = () => {
  if (!form.value.name?.trim()) return false
  if (form.value.type === 'local') return (form.value.commandStr ?? '').trim().length > 0
  return (form.value.url ?? '').trim().length > 0
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
      <t-button theme="primary" @click="openAddDialog">
        <template #icon><AddIcon /></template>
        {{ t('mcp.addServer') }}
      </t-button>
    </div>

    <t-loading :loading="loading">
      <div v-if="servers.length === 0 && !loading" class="text-center py-12">
        <div class="flex flex-col items-center gap-4">
          <ServerIcon size="48" class="text-text-muted" />
          <h4 class="text-lg font-bold">{{ t('mcp.noServers') }}</h4>
          <p class="text-text-secondary">{{ t('mcp.noServersDesc') }}</p>
          <t-button @click="openAddDialog">{{ t('mcp.addFirst') }}</t-button>
        </div>
      </div>

      <t-list v-else :split="true" class="bg-bg-card rounded-lg border border-border">
        <t-list-item v-for="item in servers" :key="item.name" class="hover:bg-bg-hover transition-colors">
          <div class="flex items-center justify-between w-full">
            <div class="flex items-center gap-3">
              <div class="p-2 bg-bg-secondary rounded">
                <ServerIcon />
              </div>
              <div class="flex flex-col">
                <div class="flex items-center gap-2 flex-wrap">
                  <span class="font-medium text-text-primary">{{ item.name }}</span>
                  <t-tag size="small" theme="default">{{ item.type === 'local' ? t('mcp.typeLocal') : t('mcp.typeRemote') }}</t-tag>
                  <t-tag :theme="statusTheme(item.status)" size="small">
                    <template v-if="checkingName === item.name">
                      <LoadingIcon class="animate-spin inline" />
                      {{ t('mcp.checking') }}
                    </template>
                    <template v-else>
                      {{ statusLabel(item.status) }}
                    </template>
                  </t-tag>
                  <t-tag :theme="item.enabled ? 'success' : 'default'" size="small">
                    {{ item.enabled ? t('common.enabled') : t('common.disabled') }}
                  </t-tag>
                </div>
                <span class="text-xs text-text-secondary">{{ displaySubtitle(item) }}</span>
                <span v-if="item.error" class="text-xs text-red-500">{{ item.error }}</span>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <t-button
                variant="text"
                shape="square"
                :theme="item.enabled ? 'success' : 'default'"
                :disabled="checkingName === item.name"
                @click="toggleServer(item)"
              >
                <template #icon>
                  <PlayCircleIcon v-if="item.enabled" />
                  <StopCircleIcon v-else />
                </template>
              </t-button>
              <t-button variant="text" shape="square" @click="openEditDialog(item)">
                <template #icon><EditIcon /></template>
              </t-button>
              <t-button variant="text" shape="square" theme="danger" @click="deleteServer(item)">
                <template #icon><DeleteIcon /></template>
              </t-button>
            </div>
          </div>
        </t-list-item>
      </t-list>
    </t-loading>

    <t-dialog
      v-model:visible="showDialog"
      :header="editingServer ? t('mcp.editServer') : t('mcp.addServer')"
      :confirm-btn="t('common.save')"
      :cancel-btn="t('common.cancel')"
      @confirm="saveServer"
      @cancel="showDialog = false"
      :confirm-loading="saving"
      :confirm-btn-props="{ disabled: !canSave() }"
    >
      <div class="space-y-4">
        <div class="grid gap-2">
          <label class="text-sm font-medium">{{ t('mcp.name') }}</label>
          <t-input
            v-model="form.name"
            :placeholder="t('mcp.namePlaceholder')"
            :disabled="!!editingServer"
            autofocus
          />
        </div>
        <div class="grid gap-2">
          <label class="text-sm font-medium">{{ t('mcp.type') }}</label>
          <t-radio-group v-model="form.type" :disabled="!!editingServer">
            <t-radio value="remote">{{ t('mcp.typeRemote') }}</t-radio>
            <t-radio value="local">{{ t('mcp.typeLocal') }}</t-radio>
          </t-radio-group>
        </div>

        <template v-if="form.type === 'local'">
          <div class="grid gap-2">
            <label class="text-sm font-medium">{{ t('mcp.command') }}</label>
            <t-input
              v-model="form.commandStr"
              :placeholder="t('mcp.commandPlaceholder')"
            />
          </div>
        </template>

        <template v-else>
          <div class="grid gap-2">
            <label class="text-sm font-medium">{{ t('mcp.baseUrl') }}</label>
            <t-input v-model="form.url" :placeholder="t('mcp.urlPlaceholder')" />
          </div>
          <div class="grid gap-2">
            <label class="text-sm font-medium">{{ t('mcp.apiKey') }}</label>
            <t-input v-model="form.api_key" type="password" :placeholder="t('mcp.apiKeyPlaceholder')" />
          </div>
        </template>
      </div>
    </t-dialog>
  </div>
</template>
