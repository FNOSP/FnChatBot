<script setup lang="ts">
import { ref, onMounted, nextTick, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { SendIcon, MoonIcon, SunnyIcon, EnterIcon, CloseIcon } from 'tdesign-icons-vue-next'
import MessageItem from '../components/chat/MessageItem.vue'
import TaskPanel from '../components/chat/TaskPanel.vue'
import { useChatStore } from '../store/chat'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { useTheme } from '../composables/useTheme'

const { t } = useI18n()
const { toggleTheme, isDark } = useTheme()
const route = useRoute()
const chatStore = useChatStore()
const { messages, isThinking, models, currentModelId, currentModel } = storeToRefs(chatStore)

const conversationId = ref(route.params.id as string || 'default')
const inputValue = ref('')
const messagesContainer = ref<HTMLElement | null>(null)

type AttachmentStatus = 'success' | 'progress' | 'fail' | 'default'

// Basic attachment item type for the chat sender
interface AttachmentItem {
  key?: string | number
  name: string
  size?: number
  status?: AttachmentStatus
  description?: string
  url?: string
}

// Attachment list shown in the chat sender
const filesList = ref<AttachmentItem[]>([])
// Whether to show the reference header above the sender
const showReferenceHeader = ref(true)
// Tooltip visibility for model select
const allowToolTip = ref(false)

// Map model list into select options for the dropdown
const modelOptions = computed(() => {
  return models.value.map(m => ({
    value: m.id,
    label: m.name || m.model,
    disabled: false,
  }))
})

// Keep the message list scrolled to the latest content
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

onMounted(() => {
  chatStore.fetchModels()
  chatStore.connect(conversationId.value)
})

watch(messages.value, () => {
  scrollToBottom()
}, { deep: true })

// Send current input as a chat message
const handleSend = () => {
  if (!inputValue.value.trim() || isThinking.value) return

  chatStore.sendMessage(inputValue.value)
  inputValue.value = ''
}

// Change the current chat model
const handleModelChange = (value: number) => {
  chatStore.setCurrentModel(value)
}

// Handle closing the reference header above the sender
const onRemoveRef = () => {
  showReferenceHeader.value = false
}

// Handle new file selection and simulate upload progress
const handleUploadFile = ({ files }: { files: File[] }) => {
  const [file] = files
  if (!file) return

  const key = `${Date.now()}-${file.name}`

  const newFile: AttachmentItem = {
    key,
    name: file.name,
    size: file.size,
    status: 'progress',
    description: '上传中',
  }

  filesList.value = [newFile, ...filesList.value]

  setTimeout(() => {
    filesList.value = filesList.value.map(item =>
      item.key === key
        ? {
            ...item,
            status: 'success',
            description: `${Math.floor(file.size / 1024)}KB`,
          }
        : item,
    )
  }, 1000)
}

// Remove attachment from the list
const handleRemoveFile = (e: CustomEvent<AttachmentItem>) => {
  const { key } = e.detail
  filesList.value = filesList.value.filter(item => item.key !== key)
}

// Handle clicking on an attachment (placeholder for future preview/download)
const handleFileClick = (e: CustomEvent<AttachmentItem>) => {
  console.log('fileClick', e.detail)
}
</script>

<template>
  <div class="h-full flex flex-col bg-bg-card rounded-lg">
    <!-- Chat Header: align with TDesign ChatEngine style -->
    <header class="h-16 border-b border-border flex items-center justify-between px-6">
      <div class="flex items-center gap-3">
        <div class="flex flex-col">
          <span class="text-sm text-text-secondary uppercase tracking-wide">
            {{ t('chat.conversation') || 'Conversation' }}
          </span>
          <span class="text-lg font-semibold text-text-primary">
            FnChatBot
          </span>
        </div>
        <t-tag size="small" shape="round" theme="primary" variant="light-outline">
          {{ currentModel?.name || currentModel?.model || t('chat.currentModel') }}
        </t-tag>
      </div>

      <div class="flex items-center gap-4">
        <t-button
          variant="text"
          shape="square"
          @click="toggleTheme"
          aria-label="Toggle theme"
        >
          <template #icon>
            <MoonIcon v-if="isDark" />
            <SunnyIcon v-else />
          </template>
        </t-button>
      </div>
    </header>

    <!-- Main Content: messages + task sidebar -->
    <main class="flex-1 flex min-h-0 overflow-hidden">
      <section class="flex-1 flex flex-col min-w-0">
        <div
          ref="messagesContainer"
          class="flex-1 overflow-y-auto px-6 py-4 scroll-smooth"
        >
          <div
            class="max-w-4xl mx-auto space-y-4"
          >
            <div
              v-if="messages.length === 0"
              class="h-full flex flex-col items-center justify-center text-text-muted opacity-70 py-16"
            >
              <div class="text-4xl font-bold mb-3 text-gradient">FnChatBot</div>
              <p class="text-sm">
                {{ t('chat.emptyHint') || 'Ask a question or describe a coding task to get started.' }}
              </p>
            </div>

            <MessageItem 
              v-for="(msg, index) in messages" 
              :key="index" 
              :message="msg" 
            />

            <div
              v-if="isThinking && messages[messages.length-1]?.role !== 'assistant'"
              class="flex items-center gap-3 text-sm text-text-secondary px-3 py-2"
            >
              <div class="w-7 h-7 rounded-full bg-success flex items-center justify-center shrink-0">
                <span class="text-white text-[10px] font-bold">AI</span>
              </div>
              <span>{{ t('chat.thinking') || 'Thinking...' }}</span>
            </div>
          </div>
        </div>

        <!-- Input Area -->
        <footer class="border-t border-border px-6 py-3 bg-bg-card/80 backdrop-blur">
          <div class="max-w-4xl mx-auto">
            <t-chat-sender
              v-model="inputValue"
              class="chat-sender"
              :textarea-props="{
                placeholder: t('chat.placeholder') || 'Send a message...',
              }"
              :attachments-props="{
                items: filesList,
                overflow: 'scrollX',
              }"
              :loading="isThinking"
              @send="handleSend"
              @file-select="handleUploadFile"
              @file-click="handleFileClick"
              @remove="handleRemoveFile"
            >
              <template #suffix="{ renderPresets }">
                <component :is="renderPresets([{ name: 'uploadImage' }, { name: 'uploadAttachment' }])" />
              </template>
              <template #footer-prefix>
                <div class="model-select">
                  <t-tooltip
                    v-model:visible="allowToolTip"
                    :content="t('chat.switchModel') || '切换模型'"
                    trigger="hover"
                  >
                    <t-select
                      :value="currentModelId"
                      :options="modelOptions"
                      @change="handleModelChange"
                      @focus="allowToolTip = false"
                      :placeholder="t('chat.selectModel')"
                      size="small"
                      style="min-width: 220px;"
                    />
                  </t-tooltip>
                  <div class="sender-helper">
                    <div class="sender-helper-main">
                      <EnterIcon
                        :size="'18px'"
                        :style="{ color: 'var(--td-text-color-disabled)', transform: 'scaleX(-1)', padding: '2px' }"
                      />
                      <span>
                        {{ t('chat.helper') || 'Press Enter to send, Shift+Enter for new line.' }}
                      </span>
                    </div>
                    <span class="sender-helper-sub">
                      FnChatBot can make mistakes. Consider checking important information.
                    </span>
                  </div>
                </div>
              </template>
              <template v-if="showReferenceHeader" #inner-header>
                <div
                  :style="{
                    display: 'flex',
                    width: '100%',
                    marginBottom: '8px',
                    paddingBottom: '8px',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    borderBottom: '1px solid var(--td-component-stroke)',
                  }"
                >
                  <div :style="{ flex: 1, display: 'flex', alignItems: 'center' }">
                    <EnterIcon
                      :size="'20px'"
                      :style="{ color: 'var(--td-text-color-disabled)', transform: 'scaleX(-1)', padding: '6px' }"
                    />
                    <p :style="{ fontSize: '14px', color: 'var(--td-text-color-placeholder)', marginLeft: '4px' }">
                      “牛顿第一定律（惯性定律）仅适用于惯性参考系，而不适用于非惯性参考系。”
                    </p>
                  </div>
                  <CloseIcon
                    :size="'20px'"
                    :style="{ color: 'var(--td-text-color-disabled)', padding: '6px', cursor: 'pointer' }"
                    @click="onRemoveRef"
                  />
                </div>
              </template>
            </t-chat-sender>
          </div>
        </footer>
      </section>

      <!-- Task Panel Sidebar (Right) -->
      <aside
        v-if="chatStore.currentTasks.length > 0"
        class="w-80 border-l border-border bg-bg-secondary/60 p-4 overflow-y-auto hidden xl:block"
      >
        <TaskPanel :tasks="chatStore.currentTasks" />
      </aside>
    </main>
  </div>
</template>

<style scoped>
.chat-sender .model-select {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.chat-sender .sender-helper {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
  font-size: 12px;
  color: var(--td-text-color-placeholder);
}

.chat-sender .sender-helper-main {
  display: flex;
  align-items: center;
  gap: 4px;
}

.chat-sender .sender-helper-sub {
  font-size: 11px;
  opacity: 0.8;
}
</style>

