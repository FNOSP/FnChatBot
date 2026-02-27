<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { http } from '../../services/http'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'
import { 
  ToolsIcon, 
  AddIcon, 
  DeleteIcon, 
  UploadIcon,
  FileIcon
} from 'tdesign-icons-vue-next'

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
const uploadFiles = ref<any[]>([]) // TDesign upload files
const uploading = ref(false)
const uploadError = ref('')

// --- Actions ---
// Load all installed skills from backend
const fetchSkills = async () => {
  loading.value = true
  try {
    const res = await http.get('/skills')
    skills.value = res.data
  } catch (e) {
    console.error('Failed to fetch skills', e)
    MessagePlugin.error(t('common.error'))
  } finally {
    loading.value = false
  }
}

// Toggle single skill enabled flag with optimistic UI update
const toggleSkill = async (skill: Skill) => {
  // Optimistic update
  const newValue = !skill.enabled
  skill.enabled = newValue
  
  try {
    await http.patch(`/skills/${skill.id}`, { enabled: newValue })
    MessagePlugin.success(skill.enabled ? t('common.enabled') : t('common.disabled'))
  } catch (e) {
    // Revert
    skill.enabled = !newValue
    console.error('Failed to toggle skill', e)
    MessagePlugin.error(t('common.error'))
  }
}

// Delete a single skill by id
const deleteSkill = async (id: number) => {
  const confirm = window.confirm(t('common.deleteConfirm'))
  if (!confirm) return
  try {
    await http.delete(`/skills/${id}`)
    await fetchSkills()
    MessagePlugin.success(t('common.success'))
  } catch (e) {
    console.error('Failed to delete skill', e)
    MessagePlugin.error(t('common.error'))
  }
}

// Capture selected file and prevent auto upload
const beforeUpload = (file: File) => {
    // TDesign upload might handle file list, but we want single file control
    uploadError.value = ''
    return false // Return false to prevent auto upload
}

// Upload selected skill file to backend
const uploadSkill = async () => {
  if (uploadFiles.value.length === 0) return
  
  const file = uploadFiles.value[0].raw
  if (!file) return

  uploading.value = true
  uploadError.value = ''
  
  const formData = new FormData()
  formData.append('file', file)
  
  try {
    await http.post('/skills/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    await fetchSkills()
    showUploadDialog.value = false
    uploadFiles.value = []
    MessagePlugin.success(t('common.success'))
  } catch (e: any) {
    console.error('Failed to upload skill', e)
    uploadError.value = e.response?.data?.error || 'Upload failed'
    MessagePlugin.error(uploadError.value)
  } finally {
    uploading.value = false
  }
}

const handleRemoveFile = () => {
    uploadFiles.value = []
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
        <h3 class="flex items-center gap-2 text-lg font-bold">
            <ToolsIcon />
            {{ t('skills.title') }}
        </h3>
        <span class="text-text-secondary">
            {{ t('skills.description') }}
        </span>
      </div>
      <t-button 
        theme="primary" 
        @click="showUploadDialog = true"
      >
        <template #icon><AddIcon /></template>
        {{ t('skills.addSkill') }}
      </t-button>
    </div>

    <!-- List -->
    <t-loading :loading="loading">
        <div v-if="skills.length === 0 && !loading" class="text-center py-12">
            <div class="flex flex-col items-center gap-4">
                <ToolsIcon size="48" class="text-text-muted" />
                <h4 class="text-lg font-bold">{{ t('skills.noSkills') }}</h4>
                <p class="text-text-secondary">{{ t('skills.noSkillsDesc') }}</p>
                <t-button @click="showUploadDialog = true">{{ t('skills.uploadSkill') }}</t-button>
            </div>
        </div>

        <t-list
            v-else
            :split="true"
            class="bg-bg-card rounded-lg border border-border"
        >
            <t-list-item v-for="item in skills" :key="item.id" class="hover:bg-bg-hover transition-colors">
                <div class="flex items-center justify-between w-full">
                    <div class="flex items-center gap-3">
                        <div class="p-2 bg-bg-secondary rounded">
                            <ToolsIcon />
                        </div>
                        <div class="flex flex-col">
                            <div class="flex items-center gap-2">
                                <span class="font-medium text-text-primary">{{ item.name }}</span>
                                <t-tag :theme="item.enabled ? 'success' : 'default'" size="small">
                                    {{ item.enabled ? t('common.enabled') : t('common.disabled') }}
                                </t-tag>
                            </div>
                            <span class="text-xs text-text-secondary">{{ item.description }}</span>
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <t-switch 
                            :value="item.enabled" 
                            @change="() => toggleSkill(item)"
                            size="small"
                        />
                        <t-button
                            variant="text"
                            shape="square"
                            theme="danger"
                            @click="deleteSkill(item.id)"
                        >
                            <template #icon><DeleteIcon /></template>
                        </t-button>
                    </div>
                </div>
            </t-list-item>
        </t-list>
    </t-loading>

    <!-- Upload Dialog -->
    <t-dialog
        v-model:visible="showUploadDialog"
        :header="t('skills.uploadTitle')"
        :confirm-btn="t('skills.uploadSkill')"
        :cancel-btn="t('common.cancel')"
        @confirm="uploadSkill"
        @cancel="showUploadDialog = false"
        :confirm-loading="uploading"
        :confirm-btn-props="{ disabled: uploadFiles.length === 0 }"
    >
        <t-upload
            v-model="uploadFiles"
            action="#"
            :before-upload="beforeUpload"
            :auto-upload="false"
            theme="custom"
            draggable
            accept=".md,.zip,.json"
            class="w-full"
        >
             <div class="p-8 text-center cursor-pointer border-2 border-dashed border-border rounded-lg hover:border-brand transition-colors bg-bg-secondary/30">
                <UploadIcon size="32" class="text-text-muted mb-4 mx-auto" />
                <div class="mt-2">
                    <div v-if="uploadFiles.length > 0" class="flex flex-col items-center">
                        <div class="flex items-center gap-2">
                            <FileIcon /> 
                            <span>{{ t('skills.fileSelected', { name: uploadFiles[0].name }) }}</span>
                            <t-button 
                                variant="text" 
                                shape="square" 
                                theme="danger"
                                size="small"
                                @click.stop="handleRemoveFile"
                            >
                                <template #icon><DeleteIcon /></template>
                            </t-button>
                        </div>
                    </div>
                    <div v-else class="text-text-primary font-medium">
                        {{ t('skills.dragDrop') }}
                    </div>
                </div>
                <div class="text-xs text-text-secondary mt-2">
                    {{ t('skills.supports') }}
                </div>
            </div>
        </t-upload>
        <div v-if="uploadError" class="text-sm text-error mt-2">
            {{ uploadError }}
        </div>
    </t-dialog>
  </div>
</template>
