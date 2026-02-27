<script setup lang="ts">
import { ref, onMounted, nextTick, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { SendIcon, AddIcon, MoonIcon, SunnyIcon } from 'tdesign-icons-vue-next'
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

const modelOptions = computed(() => {
  return models.value.map(m => ({
    value: m.id,
    label: m.name || m.model,
    disabled: false
  }))
})

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

const handleSend = () => {
  if (!inputValue.value.trim() || isThinking.value) return
  
  chatStore.sendMessage(inputValue.value)
  inputValue.value = ''
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

const handleModelChange = (value: number) => {
  chatStore.setCurrentModel(value)
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
        <t-select
          :value="currentModelId"
          :options="modelOptions"
          @change="handleModelChange"
          :placeholder="t('chat.selectModel')"
          style="min-width: 220px;"
          size="medium"
          borderless
        />
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
            <div class="relative flex items-end gap-3">
              <t-textarea 
                v-model="inputValue"
                @keydown="handleKeydown"
                :placeholder="t('chat.placeholder') || 'Send a message...'" 
                :autosize="{ minRows: 1, maxRows: 5 }"
                class="w-full"
              />
              <t-button 
                theme="primary"
                shape="circle"
                variant="base"
                size="large"
                class="mb-1"
                @click="handleSend"
                :disabled="!inputValue.trim() || isThinking"
              >
                <template #icon><SendIcon /></template>
              </t-button>
            </div>
            <div class="flex items-center justify-between mt-2 text-xs text-text-muted">
              <span>
                {{ t('chat.helper') || 'Press Enter to send, Shift+Enter for new line.' }}
              </span>
              <span>
                FnChatBot can make mistakes. Consider checking important information.
              </span>
            </div>
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

