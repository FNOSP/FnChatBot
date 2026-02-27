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
  <div class="flex h-full overflow-hidden bg-bg-primary">
    <main class="flex-1 flex flex-col min-w-0 overflow-hidden relative">
      <!-- Top Bar -->
      <div class="h-14 border-b border-border flex items-center px-4 justify-between bg-bg-card">
        <div class="flex items-center gap-3">
          <t-select
            :value="currentModelId"
            :options="modelOptions"
            @change="handleModelChange"
            :placeholder="t('chat.selectModel')"
            style="min-width: 200px;"
            size="medium"
            borderless
          />
        </div>
        <div class="flex items-center gap-3">
          <div class="text-xs text-text-secondary">{{ t('chat.currentModel') }}</div>
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
      </div>

      <!-- Messages Area -->
      <div class="flex flex-1 overflow-hidden relative">
        <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-6 scroll-smooth">
          <div v-if="messages.length === 0" class="h-full flex flex-col items-center justify-center text-text-muted opacity-50">
            <div class="text-4xl font-bold mb-4">FnChatBot</div>
            <p>How can I help you today?</p>
          </div>
          
          <MessageItem 
            v-for="(msg, index) in messages" 
            :key="index" 
            :message="msg" 
          />
          
          <div v-if="isThinking && messages[messages.length-1]?.role !== 'assistant'" class="flex gap-4 p-4">
             <div class="w-8 h-8 rounded-full bg-success flex items-center justify-center shrink-0">
               <span class="text-white text-xs font-bold">AI</span>
             </div>
             <div class="flex items-center text-sm text-text-secondary">
               Thinking...
             </div>
          </div>
        </div>

        <!-- Task Panel Sidebar (Right) -->
        <div v-if="chatStore.currentTasks.length > 0" class="w-80 border-l border-border bg-bg-card/50 p-4 overflow-y-auto hidden md:block">
          <TaskPanel :tasks="chatStore.currentTasks" />
        </div>
      </div>

      <!-- Input Area -->
      <div class="p-4 border-t border-border bg-bg-card">
        <div class="max-w-4xl mx-auto relative">
          <t-textarea 
            v-model="inputValue"
            @keydown="handleKeydown"
            :placeholder="t('chat.placeholder') || 'Send a message...'" 
            :autosize="{ minRows: 1, maxRows: 5 }"
            class="w-full"
          />
          <div class="absolute right-3 bottom-3">
             <t-button 
              theme="primary"
              shape="square"
              variant="text"
              @click="handleSend"
              :disabled="!inputValue.trim() || isThinking"
            >
              <template #icon><SendIcon /></template>
            </t-button>
          </div>
        </div>
        <div class="text-center text-xs text-text-muted mt-2">
          FnChatBot can make mistakes. Consider checking important information.
        </div>
      </div>
    </main>
  </div>
</template>

