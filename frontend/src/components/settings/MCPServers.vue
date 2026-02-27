<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { http } from '../../services/http'
import { useI18n } from 'vue-i18n'
import { 
  List, 
  ListItem,
  Button, 
  Modal, 
  Input, 
  Typography, 
  Spin,
  Tag,
  Space,
  Empty
} from '@kousum/semi-ui-vue'
import { 
  IconPlus, 
  IconEdit, 
  IconDelete, 
  IconServer,
  IconPlayCircle,
  IconStop
} from '@kousum/semi-icons-vue'

const { Title, Text } = Typography
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
  if (!confirm(t('common.deleteConfirm'))) return
  try {
    await http.delete(`/mcp/${id}`)
    await fetchServers()
  } catch (e) {
    console.error('Failed to delete server', e)
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
  } catch (e) {
    console.error('Failed to save server', e)
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
    } catch (e) {
        // Revert on error
        server.enabled = !newValue
        console.error('Failed to toggle server', e)
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
        <Title :heading="3" class="flex items-center gap-2">
            <IconServer />
            {{ t('mcp.title') }}
        </Title>
        <Text type="secondary">
            {{ t('mcp.description') }}
        </Text>
      </div>
      <Button 
        theme="solid" 
        type="primary" 
        :icon="IconPlus" 
        @click="openAddDialog"
      >
        {{ t('mcp.addServer') }}
      </Button>
    </div>

    <!-- List -->
    <Spin :spinning="loading">
        <div v-if="servers.length === 0 && !loading" class="text-center py-12">
            <Empty 
                :image="IconServer" 
                :title="t('mcp.noServers')" 
                :description="t('mcp.noServersDesc')"
            >
                 <Button @click="openAddDialog">{{ t('mcp.addFirst') }}</Button>
            </Empty>
        </div>

        <List
            v-else
            bordered
            class="bg-card rounded-lg"
        >
                <ListItem v-for="item in servers" :key="item.id" class="hover:bg-accent/50 transition-colors">
                    <div class="flex items-center gap-3">
                        <div class="p-2 bg-muted rounded">
                            <IconServer />
                        </div>
                        <div class="flex flex-col">
                            <div class="flex items-center gap-2">
                                <span class="font-medium">{{ item.name }}</span>
                                <Tag :color="item.enabled ? 'green' : 'grey'" size="small">
                                    {{ item.enabled ? t('common.enabled') : t('common.disabled') }}
                                </Tag>
                            </div>
                            <Text type="secondary">{{ item.base_url }}</Text>
                        </div>
                    </div>
                    <template #extra>
                        <Space>
                            <Button
                                :icon="item.enabled ? IconPlayCircle : IconStop"
                                theme="borderless"
                                type="tertiary"
                                :style="{ color: item.enabled ? 'var(--semi-color-success)' : 'var(--semi-color-text-2)' }"
                                @click="toggleServer(item)"
                            />
                            <Button
                                :icon="IconEdit"
                                theme="borderless"
                                type="tertiary"
                                @click="openEditDialog(item)"
                            />
                            <Button
                                :icon="IconDelete"
                                theme="borderless"
                                type="danger"
                                @click="deleteServer(item.id)"
                            />
                        </Space>
                    </template>
                </ListItem>
        </List>
    </Spin>

    <!-- Dialog -->
    <Modal
        :visible="showDialog"
        :title="editingServer ? t('mcp.editServer') : t('mcp.addServer')"
        :okText="t('common.save')"
        :cancelText="t('common.cancel')"
        @ok="saveServer"
        @cancel="showDialog = false"
        :confirmLoading="saving"
        :okButtonProps="{ disabled: !form.name || !form.base_url }"
    >
        <div class="space-y-4">
            <div class="grid gap-2">
                <label class="text-sm font-medium">{{ t('mcp.name') }}</label>
                <Input 
                    :value="form.name"
                    @change="val => (form.name = val)"
                    :placeholder="t('mcp.namePlaceholder')"
                    autofocus
                />
            </div>
            <div class="grid gap-2">
                <label class="text-sm font-medium">{{ t('mcp.baseUrl') }}</label>
                <Input 
                    :value="form.base_url"
                    @change="val => (form.base_url = val)"
                    :placeholder="t('mcp.urlPlaceholder')"
                />
            </div>
            <div class="grid gap-2">
                <label class="text-sm font-medium">{{ t('mcp.apiKey') }}</label>
                <Input 
                    :value="form.api_key"
                    @change="val => (form.api_key = val)"
                    mode="password"
                    :placeholder="t('mcp.apiKeyPlaceholder')"
                />
            </div>
        </div>
    </Modal>
  </div>
</template>
