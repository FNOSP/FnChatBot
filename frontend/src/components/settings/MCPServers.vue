<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { http } from '../../services/http'
import { useI18n } from 'vue-i18n'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import {
  AddIcon,
  EditIcon,
  DeleteIcon,
  ServerIcon,
  LoadingIcon,
  ChevronRightIcon,
  SettingIcon,
  RefreshIcon,
  CheckIcon,
  ErrorCircleIcon,
  HelpCircleIcon,
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

// Form for edit dialog: command as string, env as pairs
const defaultForm = (): Partial<MCPServerInfo> & { name: string; type: MCPType; enabled: boolean; commandStr?: string; envPairs?: { key: string; value: string }[] } => ({
  name: '',
  type: 'remote',
  url: '',
  api_key: '',
  enabled: true,
  status: 'unknown',
  commandStr: '',
  envPairs: [],
})

// Avatar color palette (deterministic from name hash)
const AVATAR_COLORS = [
  'rgb(59, 130, 246)',   // blue
  'rgb(34, 197, 94)',    // green
  'rgb(239, 68, 68)',    // red
  'rgb(168, 85, 247)',   // purple
  'rgb(107, 114, 128)',  // gray
  'rgb(234, 88, 12)',    // orange
  'rgb(14, 165, 233)',   // sky
  'rgb(99, 102, 241)',   // indigo
]

function avatarColor(name: string): string {
  let hash = 0
  for (let i = 0; i < name.length; i++) hash = ((hash << 5) - hash) + name.charCodeAt(i)
  return AVATAR_COLORS[Math.abs(hash) % AVATAR_COLORS.length]
}

function avatarLetter(name: string): string {
  return (name?.charAt(0) ?? '?').toUpperCase()
}

// --- State ---
const servers = ref<MCPServerInfo[]>([])
const loading = ref(false)
const showEditDialog = ref(false)
const showJsonDialog = ref(false)
const showLogsDialog = ref(false)
const editingServer = ref<MCPServerInfo | null>(null)
const form = ref(defaultForm())
const saving = ref(false)
const checkingName = ref<string | null>(null)
const restartingName = ref<string | null>(null)
const expandedName = ref<string | null>(null)
const jsonRaw = ref('')
const jsonSaving = ref(false)
const logsServer = ref<MCPServerInfo | null>(null)

// --- Fetch ---
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

// --- Add: Configure Manually (JSON) ---
const openJsonDialog = () => {
  jsonRaw.value = ''
  showJsonDialog.value = true
}

// Parse standard Claude Desktop mcpServers JSON and convert to backend payloads
function parseMcpServersJson(raw: string): { name: string; payload: Record<string, unknown> }[] {
  // Strip single-line comments for lenient parsing
  const cleaned = raw.replace(/\/\/[^\n]*/g, '').trim()
  const parsed = JSON.parse(cleaned) as Record<string, unknown>
  const mcpServers = parsed.mcpServers as Record<string, Record<string, unknown>> | undefined
  if (!mcpServers || typeof mcpServers !== 'object') {
    throw new Error(t('mcp.jsonParseError'))
  }
  const result: { name: string; payload: Record<string, unknown> }[] = []
  for (const [name, config] of Object.entries(mcpServers)) {
    if (!config || typeof config !== 'object') continue
    const command = config.command as string | undefined
    const args = config.args as string[] | undefined
    const url = config.url as string | undefined
    const env = config.env as Record<string, string> | undefined
    const enabled = config.enabled !== false
    if (url) {
      result.push({
        name,
        payload: {
          name,
          type: 'remote',
          url,
          enabled,
          api_key: config.api_key ?? '',
          headers: config.headers ?? undefined,
        },
      })
    } else if (command !== undefined) {
      const cmdArr = Array.isArray(args) ? [command, ...args] : [command]
      result.push({
        name,
        payload: {
          name,
          type: 'local',
          command: cmdArr,
          enabled,
          env: env ?? undefined,
        },
      })
    }
  }
  if (result.length === 0) throw new Error(t('mcp.jsonParseError'))
  return result
}

const saveFromJson = async () => {
  try {
    const items = parseMcpServersJson(jsonRaw.value)
    jsonSaving.value = true
    const addedNames: string[] = []
    for (const { name, payload } of items) {
      await http.post('/mcp', payload)
      addedNames.push(name)
    }
    await fetchServers()
    showJsonDialog.value = false
    MessagePlugin.success(t('common.success'))
    // Check status for each newly added server (if enabled)
    for (const name of addedNames) {
      const server = servers.value.find((s) => s.name === name)
      if (server?.enabled) {
        checkingName.value = name
        try {
          await http.post(`/mcp/${encodeURIComponent(name)}/check`)
        } finally {
          checkingName.value = null
        }
        await fetchServers()
      }
    }
  } catch (e) {
    if (e instanceof SyntaxError) {
      MessagePlugin.error(t('mcp.jsonParseError'))
    } else {
      MessagePlugin.error((e as Error).message || t('common.error'))
    }
  } finally {
    jsonSaving.value = false
  }
}

// --- Edit dialog ---
const openEditDialog = (server: MCPServerInfo) => {
  editingServer.value = server
  const commandStr = Array.isArray(server.command) ? server.command.join(' ') : ''
  const envPairs = Object.entries(server.env ?? {}).map(([key, value]) => ({ key, value }))
  form.value = {
    name: server.name,
    type: server.type,
    url: server.url ?? '',
    api_key: server.api_key ?? '',
    enabled: server.enabled,
    status: server.status,
    commandStr,
    envPairs: envPairs.length ? envPairs : [],
  }
  showEditDialog.value = true
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
    const env: Record<string, string> = {}
    for (const p of form.value.envPairs ?? []) {
      if (p.key.trim()) env[p.key.trim()] = p.value
    }
    if (Object.keys(env).length) payload.env = env
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
    showEditDialog.value = false
    MessagePlugin.success(t('common.success'))
  } catch (e) {
    console.error('Failed to save server', e)
    MessagePlugin.error(t('common.error'))
  } finally {
    saving.value = false
  }
}

const canSave = () => {
  if (!form.value.name?.trim()) return false
  if (form.value.type === 'local') return (form.value.commandStr ?? '').trim().length > 0
  return (form.value.url ?? '').trim().length > 0
}

// --- Delete with TDesign confirm modal ---
const deleteServer = (server: MCPServerInfo) => {
  const dialog = DialogPlugin.confirm({
    header: t('mcp.deleteConfirmTitle'),
    body: t('common.deleteConfirm'),
    confirmBtn: t('common.confirm'),
    cancelBtn: t('common.cancel'),
    theme: 'warning',
    onConfirm: async () => {
      try {
        await http.delete(`/mcp/${encodeURIComponent(server.name)}`)
        await fetchServers()
        MessagePlugin.success(t('common.success'))
        dialog.hide()
        dialog.destroy()
      } catch (e) {
        console.error('Failed to delete server', e)
        MessagePlugin.error(t('common.error'))
      }
    },
  })
}

// --- Toggle enable/disable (right-side switch) ---
// Use newValue from Switch @change so UI state stays in sync; always refetch after success.
const toggleServer = async (server: MCPServerInfo, newValue?: boolean) => {
  const targetEnabled = newValue !== undefined ? newValue : !server.enabled
  try {
    const payload: Record<string, unknown> = {
      type: server.type,
      enabled: targetEnabled,
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
    MessagePlugin.success(targetEnabled ? t('common.enabled') : t('common.disabled'))
    if (targetEnabled) {
      checkingName.value = server.name
      try {
        await http.post(`/mcp/${encodeURIComponent(server.name)}/check`)
      } finally {
        checkingName.value = null
      }
    }
    // Always refetch so switch state matches server (fixes controlled Switch display)
    await fetchServers()
  } catch (e) {
    console.error('Failed to toggle server', e)
    MessagePlugin.error(t('common.error'))
    await fetchServers()
  }
}

// --- Restart (reuse check endpoint: disconnect + reconnect) ---
const restartServer = async (server: MCPServerInfo) => {
  restartingName.value = server.name
  try {
    await http.post(`/mcp/${encodeURIComponent(server.name)}/check`)
    await fetchServers()
    MessagePlugin.success(t('common.success'))
  } catch (e) {
    console.error('Failed to restart server', e)
    MessagePlugin.error(t('common.error'))
  } finally {
    restartingName.value = null
  }
}

// --- Logs dialog ---
const openLogsDialog = (server: MCPServerInfo) => {
  logsServer.value = server
  showLogsDialog.value = true
}

// --- Helpers ---
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

const toggleExpand = (name: string) => {
  expandedName.value = expandedName.value === name ? null : name
}

onMounted(() => {
  fetchServers()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header: title + help + refresh + Add button (manual config only) -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <h3 class="text-lg font-bold">{{ t('mcp.title') }}</h3>
        <t-tooltip :content="t('mcp.description')">
          <t-button variant="text" shape="square" theme="default">
            <template #icon><HelpCircleIcon /></template>
          </t-button>
        </t-tooltip>
      </div>
      <div class="flex items-center gap-2">
        <t-button variant="text" shape="square" :loading="loading" @click="fetchServers">
          <template #icon><RefreshIcon /></template>
        </t-button>
        <t-button theme="default" variant="outline" @click="openJsonDialog">
          <template #icon><AddIcon /></template>
          {{ t('mcp.addServer') }}
        </t-button>
      </div>
    </div>

    <t-loading :loading="loading">
      <div v-if="servers.length === 0 && !loading" class="text-center py-12">
        <div class="flex flex-col items-center gap-4">
          <ServerIcon size="48" class="text-text-muted" />
          <h4 class="text-lg font-bold">{{ t('mcp.noServers') }}</h4>
          <p class="text-text-secondary">{{ t('mcp.noServersDesc') }}</p>
          <t-button @click="openJsonDialog">{{ t('mcp.addFirst') }}</t-button>
        </div>
      </div>

      <!-- List: chevron | avatar | name+status | gear | toggle -->
      <div v-else class="rounded-lg border border-border overflow-hidden bg-bg-card">
        <div
          v-for="item in servers"
          :key="item.name"
          class="border-b border-border last:border-b-0"
        >
          <div
            class="flex items-center gap-3 px-4 py-3 hover:bg-bg-hover transition-colors cursor-pointer"
            @click="toggleExpand(item.name)"
          >
            <div class="flex-shrink-0 w-6 flex justify-center">
              <ChevronRightIcon
                :class="{ 'rotate-90': expandedName === item.name }"
                class="transition-transform text-text-secondary"
              />
            </div>
            <div
              class="w-9 h-9 rounded flex items-center justify-center text-white font-semibold flex-shrink-0"
              :style="{ backgroundColor: avatarColor(item.name) }"
            >
              {{ avatarLetter(item.name) }}
            </div>
            <div class="flex-1 min-w-0 flex items-center gap-2">
              <span class="font-medium text-text-primary truncate">{{ item.name }}</span>
              <template v-if="checkingName === item.name || restartingName === item.name">
                <LoadingIcon class="animate-spin text-base text-brand" />
                <span class="text-sm text-text-secondary">{{ checkingName === item.name ? t('mcp.checking') : t('mcp.restarting') }}</span>
              </template>
              <template v-else>
                <CheckIcon v-if="item.status === 'connected'" class="text-green-500" size="16px" />
                <ErrorCircleIcon v-else-if="item.status === 'failed' || item.status === 'disabled'" class="text-red-500" size="16px" />
                <span v-else class="w-4 h-4 rounded-full border-2 border-amber-500"></span>
              </template>
            </div>
            <div class="flex items-center gap-1 flex-shrink-0" @click.stop>
              <t-dropdown
                :options="[
                  { content: t('common.edit'), value: 'edit' },
                  { content: t('mcp.restart'), value: 'restart' },
                  { content: t('mcp.logs'), value: 'logs' },
                  { content: t('common.delete'), value: 'delete', theme: 'error' },
                ]"
                trigger="click"
                @click="(data) => {
                  if (data.value === 'edit') openEditDialog(item)
                  else if (data.value === 'restart') restartServer(item)
                  else if (data.value === 'logs') openLogsDialog(item)
                  else if (data.value === 'delete') deleteServer(item)
                }"
              >
                <t-button variant="text" shape="square" theme="default">
                  <template #icon><SettingIcon /></template>
                </t-button>
              </t-dropdown>
            </div>
            <div class="flex-shrink-0" @click.stop>
              <t-switch
                :value="item.enabled"
                :loading="checkingName === item.name"
                :disabled="checkingName === item.name"
                @change="(val) => toggleServer(item, val)"
              />
            </div>
          </div>
          <!-- Expanded row: subtitle + error -->
          <div v-if="expandedName === item.name" class="px-4 pb-3 pt-0 pl-[3.25rem] text-sm text-text-secondary border-t border-border bg-bg-secondary/50">
            <div>{{ displaySubtitle(item) || '—' }}</div>
            <div v-if="item.error" class="text-red-500 mt-1">{{ item.error }}</div>
          </div>
        </div>
      </div>
    </t-loading>

    <!-- Edit dialog (form) -->
    <t-dialog
      v-model:visible="showEditDialog"
      :header="editingServer ? t('mcp.editServer') : t('mcp.addServer')"
      :confirm-btn="t('common.save')"
      :cancel-btn="t('common.cancel')"
      @confirm="saveServer"
      @cancel="showEditDialog = false"
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
            <t-input v-model="form.commandStr" :placeholder="t('mcp.commandPlaceholder')" />
          </div>
          <div class="grid gap-2">
            <div class="flex items-center justify-between">
              <label class="text-sm font-medium">{{ t('mcp.env') }}</label>
              <t-button size="small" variant="text" @click="(form.envPairs = form.envPairs ?? []).push({ key: '', value: '' })">
                <template #icon><AddIcon /></template>
                {{ t('mcp.envAdd') }}
              </t-button>
            </div>
            <div v-for="(pair, idx) in (form.envPairs ?? [])" :key="idx" class="flex gap-2 items-center">
              <t-input v-model="pair.key" :placeholder="t('mcp.envKey')" class="flex-1" />
              <t-input v-model="pair.value" :placeholder="t('mcp.envValue')" class="flex-1" />
              <t-button variant="text" theme="danger" shape="square" @click="(form.envPairs ?? []).splice(idx, 1)">
                <template #icon><DeleteIcon /></template>
              </t-button>
            </div>
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

    <!-- Configure Manually (JSON) dialog -->
    <t-dialog
      v-model:visible="showJsonDialog"
      :header="t('mcp.configureManually')"
      :confirm-btn="t('common.confirm')"
      :cancel-btn="t('common.cancel')"
      :confirm-loading="jsonSaving"
      @confirm="saveFromJson"
      @cancel="showJsonDialog = false"
    >
      <div class="space-y-3">
        <p class="text-sm text-text-secondary">{{ t('mcp.jsonInstruction') }}</p>
        <t-textarea
          v-model="jsonRaw"
          :placeholder="t('mcp.jsonPlaceholder')"
          class="font-mono text-sm"
          :autosize="{ minRows: 12, maxRows: 20 }"
        />
        <t-alert theme="warning" :message="t('mcp.configRiskWarning')" class="text-xs" />
      </div>
    </t-dialog>

    <!-- Logs dialog -->
    <t-dialog
      v-model:visible="showLogsDialog"
      :header="logsServer ? `${t('mcp.logs')}: ${logsServer.name}` : t('mcp.logs')"
      :footer="false"
      @close="logsServer = null"
    >
      <div class="font-mono text-sm whitespace-pre-wrap break-words">
        {{ logsServer?.error || t('mcp.noLogs') }}
      </div>
    </t-dialog>
  </div>
</template>
