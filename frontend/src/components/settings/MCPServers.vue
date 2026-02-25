<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
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
const fetchServers = async () => {
  loading.value = true
  try {
    const res = await axios.get('http://localhost:8080/api/mcp')
    servers.value = res.data
  } catch (e) {
    console.error('Failed to fetch MCP servers', e)
  } finally {
    loading.value = false
  }
}

const openAddDialog = () => {
  editingServer.value = null
  form.value = { name: '', base_url: '', api_key: '', enabled: true }
  showDialog.value = true
}

const openEditDialog = (server: MCPConfig) => {
  editingServer.value = server
  form.value = { ...server }
  showDialog.value = true
}

const deleteServer = async (id: number) => {
  if (!confirm(t('common.deleteConfirm'))) return
  try {
    await axios.delete(`http://localhost:8080/api/mcp/${id}`)
    await fetchServers()
  } catch (e) {
    console.error('Failed to delete server', e)
  }
}

const saveServer = async () => {
  saving.value = true
  try {
    if (editingServer.value && editingServer.value.id) {
      await axios.put(`http://localhost:8080/api/mcp/${editingServer.value.id}`, form.value)
    } else {
      await axios.post('http://localhost:8080/api/mcp', form.value)
    }
    await fetchServers()
    showDialog.value = false
  } catch (e) {
    console.error('Failed to save server', e)
  } finally {
    saving.value = false
  }
}

const toggleServer = async (server: MCPConfig) => {
    if (!server.id) return
    // Optimistic update
    const newValue = !server.enabled
    server.enabled = newValue
    
    try {
        await axios.put(`http://localhost:8080/api/mcp/${server.id}`, { ...server, enabled: newValue })
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
        v-model:visible="showDialog"
        :title="editingServer ? t('mcp.editServer') : t('mcp.addServer')"
        :okText="t('common.save')"
        :cancelText="t('common.cancel')"
        @ok="saveServer"
        :confirmLoading="saving"
        :okButtonProps="{ disabled: !form.name || !form.base_url }"
    >
        <div class="space-y-4">
            <div class="grid gap-2">
                <label class="text-sm font-medium">{{ t('mcp.name') }}</label>
                <Input 
                    v-model="form.name" 
                    :placeholder="t('mcp.namePlaceholder')"
                    autofocus
                />
            </div>
            <div class="grid gap-2">
                <label class="text-sm font-medium">{{ t('mcp.baseUrl') }}</label>
                <Input 
                    v-model="form.base_url" 
                    :placeholder="t('mcp.urlPlaceholder')"
                />
            </div>
            <div class="grid gap-2">
                <label class="text-sm font-medium">{{ t('mcp.apiKey') }}</label>
                <Input 
                    v-model="form.api_key" 
                    mode="password"
                    :placeholder="t('mcp.apiKeyPlaceholder')"
                />
            </div>
        </div>
    </Modal>
  </div>
</template>
