<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { http } from '../../services/http'
import { useI18n } from 'vue-i18n'
import { 
  List, 
  ListItem,
  Button, 
  Modal, 
  Switch,
  Typography, 
  Spin,
  Tag,
  Space,
  Empty,
  Upload
} from '@kousum/semi-ui-vue'
import { 
  IconWrench, 
  IconPlus, 
  IconDelete, 
  IconUpload,
  IconFile
} from '@kousum/semi-icons-vue'

const { Title, Text } = Typography
const { t } = useI18n()

// --- Types ---
interface Skill {
  id: number
  name: string
  description: string
  enabled: boolean
}

// --- State ---
const skills = ref<Skill[]>([])
const loading = ref(false)
const showUploadDialog = ref(false)
const uploadFile = ref<File | null>(null)
const uploading = ref(false)
const uploadError = ref('')

// --- Actions ---
const fetchSkills = async () => {
  loading.value = true
  try {
    const res = await http.get('/skills')
    skills.value = res.data
  } catch (e) {
    console.error('Failed to fetch skills', e)
  } finally {
    loading.value = false
  }
}

const toggleSkill = async (skill: Skill) => {
  // Optimistic update
  const newValue = !skill.enabled
  skill.enabled = newValue
  
  try {
    await http.patch(`/skills/${skill.id}`, { enabled: newValue })
  } catch (e) {
    // Revert
    skill.enabled = !newValue
    console.error('Failed to toggle skill', e)
  }
}

const deleteSkill = async (id: number) => {
  if (!confirm(t('common.deleteConfirm'))) return
  try {
    await http.delete(`/skills/${id}`)
    await fetchSkills()
  } catch (e) {
    console.error('Failed to delete skill', e)
  }
}

const beforeUpload = (file: any) => {
    // Semi UI Upload component passes a File object
    // We want to prevent default upload and handle it manually
    uploadFile.value = file
    uploadError.value = ''
    return false // Return false to prevent auto upload
}

const uploadSkill = async () => {
  if (!uploadFile.value) return
  
  uploading.value = true
  uploadError.value = ''
  
  const formData = new FormData()
  formData.append('file', uploadFile.value)
  
  try {
    await http.post('/skills/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    await fetchSkills()
    showUploadDialog.value = false
    uploadFile.value = null
  } catch (e: any) {
    console.error('Failed to upload skill', e)
    uploadError.value = e.response?.data?.error || 'Upload failed'
  } finally {
    uploading.value = false
  }
}

const handleRemoveFile = () => {
    uploadFile.value = null
    uploadError.value = ''
}

onMounted(() => {
  fetchSkills()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <Title :heading="3" class="flex items-center gap-2">
            <IconWrench />
            {{ t('skills.title') }}
        </Title>
        <Text type="secondary">
            {{ t('skills.description') }}
        </Text>
      </div>
      <Button 
        theme="solid" 
        type="primary" 
        :icon="IconPlus" 
        @click="showUploadDialog = true"
      >
        {{ t('skills.addSkill') }}
      </Button>
    </div>

    <!-- List -->
    <Spin :spinning="loading">
        <div v-if="skills.length === 0 && !loading" class="text-center py-12">
            <Empty 
                :image="IconWrench" 
                :title="t('skills.noSkills')" 
                :description="t('skills.noSkillsDesc')"
            >
                 <Button @click="showUploadDialog = true">{{ t('skills.uploadSkill') }}</Button>
            </Empty>
        </div>

        <List
            v-else
            bordered
            class="bg-card rounded-lg"
        >
                <ListItem v-for="item in skills" :key="item.id" class="hover:bg-accent/50 transition-colors">
                    <div class="flex items-center gap-3">
                        <div class="p-2 bg-muted rounded">
                            <IconWrench />
                        </div>
                        <div class="flex flex-col">
                            <div class="flex items-center gap-2">
                                <span class="font-medium">{{ item.name }}</span>
                                <Tag :color="item.enabled ? 'green' : 'grey'" size="small">
                                    {{ item.enabled ? t('common.enabled') : t('common.disabled') }}
                                </Tag>
                            </div>
                            <Text type="secondary">{{ item.description }}</Text>
                        </div>
                    </div>
                    <template #extra>
                        <Space>
                            <Switch 
                                :checked="item.enabled" 
                                @change="() => toggleSkill(item)"
                                size="small"
                            />
                            <Button
                                :icon="IconDelete"
                                theme="borderless"
                                type="danger"
                                @click="deleteSkill(item.id)"
                            />
                        </Space>
                    </template>
                </ListItem>
        </List>
    </Spin>

    <!-- Upload Dialog -->
    <Modal
        v-model:visible="showUploadDialog"
        :title="t('skills.uploadTitle')"
        :okText="t('skills.uploadSkill')"
        :cancelText="t('common.cancel')"
        @ok="uploadSkill"
        :confirmLoading="uploading"
        :okButtonProps="{ disabled: !uploadFile }"
    >
        <Upload
            action="#"
            :beforeUpload="beforeUpload"
            :showUploadList="false"
            drag
            accept=".md,.zip,.json"
        >
             <div class="p-8 text-center cursor-pointer">
                <IconUpload size="extra-large" class="text-muted-foreground mb-4" />
                <div class="mt-2">
                    <div v-if="uploadFile" class="flex flex-col items-center">
                        <Space>
                            <IconFile /> 
                            <span>{{ t('skills.fileSelected', { name: uploadFile.name }) }}</span>
                            <Button 
                                type="tertiary" 
                                theme="borderless" 
                                :icon="IconDelete" 
                                size="small"
                                @click.stop="handleRemoveFile"
                            />
                        </Space>
                    </div>
                    <div v-else>
                        {{ t('skills.dragDrop') }}
                    </div>
                </div>
                <div class="text-xs text-muted-foreground mt-2">
                    {{ t('skills.supports') }}
                </div>
            </div>
        </Upload>
        <div v-if="uploadError" class="text-sm text-red-500 mt-2">
            {{ uploadError }}
        </div>
    </Modal>
  </div>
</template>
