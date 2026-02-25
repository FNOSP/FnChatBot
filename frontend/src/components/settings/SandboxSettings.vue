<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
import {
  Layout,
  Card,
  Switch,
  Button,
  Input,
  Typography,
  Spin,
  List,
  ListItem,
  Space,
  Empty,
  Toast
} from '@kousum/semi-ui-vue'
import {
  IconPlus,
  IconDelete,
  IconShield
} from '@kousum/semi-icons-vue'

const { Title, Text } = Typography
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
    const res = await axios.get('http://localhost:8080/api/sandbox')
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
    await axios.put('http://localhost:8080/api/sandbox', { enabled: checked })
    enabled.value = checked
    Toast.success(checked ? t('common.enabled') : t('common.disabled'))
  } catch (e) {
    console.error('Failed to toggle sandbox', e)
    Toast.error(t('common.operationFailed'))
  }
}

const addPath = async () => {
  if (!newPath.value.trim()) {
    Toast.warning(t('sandbox.pathPlaceholder'))
    return
  }

  adding.value = true
  try {
    await axios.post('http://localhost:8080/api/sandbox/paths', {
      path: newPath.value.trim(),
      description: newDescription.value.trim()
    })
    paths.value.push({
      path: newPath.value.trim(),
      description: newDescription.value.trim()
    })
    newPath.value = ''
    newDescription.value = ''
    Toast.success(t('common.addSuccess'))
  } catch (e) {
    console.error('Failed to add path', e)
    Toast.error(t('common.operationFailed'))
  } finally {
    adding.value = false
  }
}

const removePath = async (path: string) => {
  try {
    await axios.delete(`http://localhost:8080/api/sandbox/paths/${encodeURIComponent(path)}`)
    paths.value = paths.value.filter(p => p.path !== path)
    Toast.success(t('common.deleteSuccess'))
  } catch (e) {
    console.error('Failed to remove path', e)
    Toast.error(t('common.operationFailed'))
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<template>
  <Layout class="h-full w-full">
    <Card :title="t('sandbox.title')" :bordered="false" class="h-full bg-transparent shadow-none">
      <template #headerExtra>
        <Text type="secondary">{{ t('sandbox.description') }}</Text>
      </template>

      <Spin :spinning="loading">
        <div class="flex flex-col gap-6">
          <div class="flex items-center justify-between p-4 border border-zinc-200 dark:border-zinc-800 rounded-lg bg-card">
            <div class="flex flex-col">
              <span class="text-base font-medium">{{ t('sandbox.enabled') }}</span>
              <span class="text-sm text-zinc-500">{{ t('sandbox.enabledDesc') }}</span>
            </div>
            <Switch
              :checked="enabled"
              @change="toggleEnabled"
            />
          </div>

          <div class="border border-zinc-200 dark:border-zinc-800 rounded-lg bg-card p-4">
            <Title :heading="5" class="mb-4 flex items-center gap-2">
              <IconShield />
              {{ t('sandbox.allowedPaths') }}
            </Title>

            <div class="flex gap-3 mb-4">
              <Input
                v-model="newPath"
                :placeholder="t('sandbox.pathPlaceholder')"
                class="flex-1"
                @keydown.enter="addPath"
              />
              <Input
                v-model="newDescription"
                :placeholder="t('sandbox.descriptionPlaceholder')"
                class="flex-1"
                @keydown.enter="addPath"
              />
              <Button
                theme="solid"
                type="primary"
                :icon="IconPlus"
                :loading="adding"
                @click="addPath"
              >
                {{ t('sandbox.addPath') }}
              </Button>
            </div>

            <div v-if="paths.length === 0 && !loading" class="text-center py-8">
              <Empty
                :image="IconShield"
                :title="t('sandbox.noPaths')"
                :description="t('sandbox.noPathsDesc')"
              />
            </div>

            <List
              v-else
              bordered
              class="bg-transparent"
            >
                <ListItem v-for="item in paths" :key="item.path" class="hover:bg-accent/50 transition-colors">
                  <div class="flex flex-col">
                    <div class="font-medium">{{ item.path }}</div>
                    <div class="text-sm text-muted-foreground">{{ item.description || '-' }}</div>
                  </div>
                  <template #extra>
                    <Button
                      :icon="IconDelete"
                      theme="borderless"
                      type="danger"
                      @click="removePath(item.path)"
                    />
                  </template>
                </ListItem>
            </List>
          </div>
        </div>
      </Spin>
    </Card>
  </Layout>
</template>

<style scoped>
:deep(.semi-card-header) {
  padding-left: 0;
  padding-right: 0;
}
:deep(.semi-card-body) {
  padding-left: 0;
  padding-right: 0;
}
</style>
