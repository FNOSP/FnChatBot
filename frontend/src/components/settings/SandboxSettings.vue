<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { http } from '../../services/http'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'
import {
  AddIcon,
  DeleteIcon,
  SecuredIcon
} from 'tdesign-icons-vue-next'

const { t } = useI18n()

interface PathItem {
  path: string
  description: string
}

interface SandboxConfig {
  enabled: boolean
  paths: PathItem[]
}

const enabled = ref(false)
const paths = ref<PathItem[]>([])
const newPath = ref('')
const newDescription = ref('')
const loading = ref(false)
const adding = ref(false)

const fetchConfig = async () => {
  loading.value = true
  try {
    const res = await http.get('/sandbox')
    enabled.value = res.data.enabled || false
    paths.value = res.data.paths || []
  } catch (e) {
    console.error('Failed to fetch sandbox config', e)
  } finally {
    loading.value = false
  }
}

const toggleEnabled = async (checked: boolean) => {
  try {
    await http.put('/sandbox', { enabled: checked })
    enabled.value = checked
    MessagePlugin.success(checked ? t('common.enabled') : t('common.disabled'))
  } catch (e) {
    console.error('Failed to toggle sandbox', e)
    MessagePlugin.error(t('common.operationFailed'))
  }
}

const addPath = async () => {
  if (!newPath.value.trim()) {
    MessagePlugin.warning(t('sandbox.pathPlaceholder'))
    return
  }

  adding.value = true
  try {
    await http.post('/sandbox/paths', {
      path: newPath.value.trim(),
      description: newDescription.value.trim()
    })
    paths.value.push({
      path: newPath.value.trim(),
      description: newDescription.value.trim()
    })
    newPath.value = ''
    newDescription.value = ''
    MessagePlugin.success(t('common.addSuccess'))
  } catch (e) {
    console.error('Failed to add path', e)
    MessagePlugin.error(t('common.operationFailed'))
  } finally {
    adding.value = false
  }
}

const removePath = async (path: string) => {
  try {
    await http.delete(`/sandbox/paths/${encodeURIComponent(path)}`)
    paths.value = paths.value.filter(p => p.path !== path)
    MessagePlugin.success(t('common.deleteSuccess'))
  } catch (e) {
    console.error('Failed to remove path', e)
    MessagePlugin.error(t('common.operationFailed'))
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<template>
  <div class="h-full w-full p-4">
    <t-card :title="t('sandbox.title')" :bordered="false" class="h-full bg-transparent shadow-none">
      <template #subtitle>
        <span class="text-text-secondary">{{ t('sandbox.description') }}</span>
      </template>

      <t-loading :loading="loading" class="w-full h-full">
        <div class="flex flex-col gap-6">
          <div class="flex items-center justify-between p-4 border border-border rounded-lg bg-bg-card">
            <div class="flex flex-col">
              <span class="text-base font-medium text-text-primary">{{ t('sandbox.enabled') }}</span>
              <span class="text-sm text-text-secondary">{{ t('sandbox.enabledDesc') }}</span>
            </div>
            <t-switch
              :value="enabled"
              @change="toggleEnabled"
            />
          </div>

          <div class="border border-border rounded-lg bg-bg-card p-4">
            <h5 class="mb-4 flex items-center gap-2 text-lg font-bold">
              <SecuredIcon />
              {{ t('sandbox.allowedPaths') }}
            </h5>

            <div class="flex gap-3 mb-4">
              <t-input
                v-model="newPath"
                :placeholder="t('sandbox.pathPlaceholder')"
                class="flex-1"
                @keydown.enter="addPath"
              />
              <t-input
                v-model="newDescription"
                :placeholder="t('sandbox.descriptionPlaceholder')"
                class="flex-1"
                @keydown.enter="addPath"
              />
              <t-button
                theme="primary"
                :loading="adding"
                @click="addPath"
              >
                <template #icon><AddIcon /></template>
                {{ t('sandbox.addPath') }}
              </t-button>
            </div>

            <div v-if="paths.length === 0 && !loading" class="text-center py-8">
              <div class="flex flex-col items-center gap-2 text-text-muted">
                <SecuredIcon size="32" />
                <span class="font-bold">{{ t('sandbox.noPaths') }}</span>
                <span class="text-sm">{{ t('sandbox.noPathsDesc') }}</span>
              </div>
            </div>

            <t-list
              v-else
              :split="true"
              class="bg-transparent"
            >
                <t-list-item v-for="item in paths" :key="item.path" class="hover:bg-bg-hover transition-colors">
                  <div class="flex flex-col">
                    <div class="font-medium text-text-primary">{{ item.path }}</div>
                    <div class="text-sm text-text-secondary">{{ item.description || '-' }}</div>
                  </div>
                  <template #action>
                    <t-button
                      theme="danger"
                      variant="text"
                      shape="square"
                      @click="removePath(item.path)"
                    >
                      <template #icon><DeleteIcon /></template>
                    </t-button>
                  </template>
                </t-list-item>
            </t-list>
          </div>
        </div>
      </t-loading>
    </t-card>
  </div>
</template>
