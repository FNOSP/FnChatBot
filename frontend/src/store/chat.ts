import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { WebSocketService } from '../services/ws'
import { http } from '../services/http'

export interface Message {
  role: 'user' | 'assistant' | 'system'
  content: string
  id?: string
  thinking?: string
  tasks?: any[]
}

export interface Model {
  id?: number
  name: string
  model: string
  provider: string
  base_url: string
  api_key: string
}

export const useChatStore = defineStore('chat', () => {
  const messages = ref<Message[]>([])
  const isThinking = ref(false)
  const currentThinking = ref('')
  const currentTasks = ref<any[]>([])
  const wsService = ref<WebSocketService | null>(null)
  const models = ref<Model[]>([])
  const currentModelId = ref<number | null>(null)

  const currentModel = computed(() => {
    if (currentModelId.value === null) return null
    return models.value.find(m => m.id === currentModelId.value) || null
  })

  const setCurrentModel = (modelId: number | null) => {
    currentModelId.value = modelId
  }

  const fetchModels = async () => {
    try {
      const res = await http.get('/models')
      models.value = res.data || []
      if (models.value.length === 0) {
        currentModelId.value = null
      }
    } catch (e) {
      console.error('Failed to fetch models', e)
      models.value = []
    }
  }

  const connect = (conversationId: string) => {
    // Determine backend URL (assuming localhost:8080 for dev)
    const backendUrl = 'ws://localhost:8080'
    wsService.value = new WebSocketService(backendUrl)
    
    wsService.value.onMessage((data) => {
      handleMessage(data)
    })
    
    wsService.value.connect(conversationId)
  }

  const sendMessage = (content: string) => {
    if (!wsService.value) return

    // Add user message locally
    messages.value.push({
      role: 'user',
      content
    })

    wsService.value.sendMessage({
      type: 'user_message',
      content
    })
    
    isThinking.value = true
    currentThinking.value = ''
    currentTasks.value = []
  }

  const handleMessage = (data: any) => {
    switch (data.type) {
      case 'thinking':
        currentThinking.value = data.content
        break
      case 'task_update':
        currentTasks.value = data.tasks
        // Update the tasks of the current assistant message if it exists
        const lastMsgForTask = messages.value[messages.value.length - 1]
        if (lastMsgForTask && lastMsgForTask.role === 'assistant') {
          lastMsgForTask.tasks = data.tasks
        }
        break
      case 'message':
        // Stream delta to the last assistant message or create new one
        const lastMsg = messages.value[messages.value.length - 1]
        if (lastMsg && lastMsg.role === 'assistant' && !data.content) { // data.content is usually empty for delta, check type
             // If we rely on delta field
             if (data.delta) {
                 lastMsg.content += data.delta
             }
             // Ensure tasks are synced
             if (currentTasks.value.length > 0) {
                 lastMsg.tasks = currentTasks.value
             }
        } else {
          // If content is present, it might be a full message or start of one
          // But our backend sends Delta for stream
          const content = data.delta || data.content || ''
          messages.value.push({
            role: 'assistant',
            content: content,
            thinking: currentThinking.value,
            tasks: currentTasks.value
          })
        }
        break
      case 'message_end':
        isThinking.value = false
        break
    }
  }

  return {
    messages,
    isThinking,
    currentThinking,
    currentTasks,
    models,
    currentModelId,
    currentModel,
    connect,
    sendMessage,
    setCurrentModel,
    fetchModels
  }
})
