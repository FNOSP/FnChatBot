<script setup lang="ts">
import { ref, onMounted, nextTick, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Send, Plus, Sun, Moon } from 'lucide-vue-next'
import { Select } from '@kousum/semi-ui-vue'
import Sidebar from '../components/layout/Sidebar.vue'
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
  <div class="flex h-screen overflow-hidden bg-background">
    <Sidebar />
    <main class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Top Bar -->
      <div class="h-12 border-b border-border flex items-center px-4 justify-between">
        <div class="flex items-center gap-3">
          <Select
            :value="currentModelId"
            :options="modelOptions"
            @change="handleModelChange"
            :placeholder="t('chat.selectModel')"
            style="min-width: 150px;"
            size="small"
          />
        </div>
        <div class="flex items-center gap-3">
          <div class="text-xs text-muted-foreground">{{ t('chat.currentModel') }}</div>
          <button
            class="p-2 rounded-md hover:bg-secondary transition-colors"
            @click="toggleTheme"
            aria-label="Toggle theme"
          >
            <Moon v-if="isDark" class="w-4 h-4" />
            <Sun v-else class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Messages Area -->
      <div class="flex flex-1 overflow-hidden">
        <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-4">
          <div v-if="messages.length === 0" class="h-full flex flex-col items-center justify-center text-muted-foreground opacity-50">
            <div class="text-4xl font-bold mb-4">FnChatBot</div>
            <p>How can I help you today?</p>
          </div>
          
          <MessageItem 
            v-for="(msg, index) in messages" 
            :key="index" 
            :message="msg" 
          />
          
          <div v-if="isThinking && messages[messages.length-1]?.role !== 'assistant'" class="flex gap-4 p-4">
             <div class="w-8 h-8 rounded-full bg-green-600 flex items-center justify-center shrink-0">
               <span class="text-white text-xs font-bold">AI</span>
             </div>
             <div class="flex items-center text-sm text-muted-foreground">
               Thinking...
             </div>
          </div>
        </div>

        <!-- Task Panel Sidebar (Right) -->
        <div v-if="chatStore.currentTasks.length > 0" class="w-80 border-l border-border bg-card/50 p-4 overflow-y-auto hidden md:block">
          <TaskPanel :tasks="chatStore.currentTasks" />
        </div>
      </div>

      <!-- Input Area -->
      <div class="p-4 border-t border-border bg-background">
        <div class="max-w-4xl mx-auto relative">
          <textarea 
            v-model="inputValue"
            @keydown="handleKeydown"
            placeholder="Send a message..." 
            class="w-full p-3 pr-12 border rounded-xl bg-input text-foreground resize-none focus:outline-none focus:ring-2 focus:ring-ring min-h-[50px] max-h-[200px]"
            rows="1"
          ></textarea>
          <button 
            @click="handleSend"
            :disabled="!inputValue.trim() || isThinking"
            class="absolute right-3 bottom-3 p-1.5 rounded-lg bg-primary text-primary-foreground hover:opacity-90 disabled:opacity-50 transition-opacity"
          >
            <Send class="w-4 h-4" />
          </button>
        </div>
        <div class="text-center text-xs text-muted-foreground mt-2">
          FnChatBot can make mistakes. Consider checking important information.
        </div>
      </div>
    </main>
  </div>
</template>
