<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { AddIcon, DeleteIcon, ChatIcon } from 'tdesign-icons-vue-next'
import { MessagePlugin } from 'tdesign-vue-next'
import { useChatStore } from '../../store/chat'
import { useTheme } from '../../composables/useTheme'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const chatStore = useChatStore()
const { isDark } = useTheme()

const currentConversationId = computed(() => {
  return route.params.id ? Number(route.params.id) : null
})

onMounted(async () => {
  await chatStore.fetchConversations()
})

const handleNewChat = async () => {
  const session = await chatStore.createConversation(t('sidebar.newChat'))
  if (session) {
    router.push(`/chat/${session.id}`)
  }
}

const handleSelectConversation = (id: number) => {
  router.push(`/chat/${id}`)
}

const handleDeleteConversation = async (e: Event, id: number) => {
  e.stopPropagation()
  const confirmed = await new Promise<boolean>((resolve) => {
    MessagePlugin.confirm?.({
      content: t('sidebar.deleteConfirm'),
      onConfirm: () => resolve(true),
      onCancel: () => resolve(false),
    })
  })
  // Fallback: use window.confirm if MessagePlugin.confirm is unavailable
  if (!MessagePlugin.confirm) {
    if (!window.confirm(t('sidebar.deleteConfirm'))) return
  } else if (!confirmed) return

  const ok = await chatStore.deleteConversation(id)
  if (ok) {
    // If we deleted the current conversation, navigate home
    if (currentConversationId.value === id) {
      router.push('/')
    }
  }
}

const formatDate = (dateStr: string) => {
  const d = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - d.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  if (diffDays === 0) return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  if (diffDays === 1) return t('common.yesterday') || 'Yesterday'
  if (diffDays < 7) return `${diffDays}${t('common.daysAgo') || 'd ago'}`
  return d.toLocaleDateString()
}
</script>

<template>
  <div
    class="conv-sidebar flex flex-col border-r border-border"
    :class="isDark ? 'bg-bg-secondary' : 'bg-bg-secondary'"
    style="width: 260px; min-width: 260px;"
  >
    <!-- Header -->
    <div class="px-4 py-4 shrink-0 border-b border-border">
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center gap-2">
          <ChatIcon class="text-brand" />
          <span class="text-sm font-semibold text-text-primary">FnChatBot</span>
        </div>
      </div>
      <!-- New Chat Button -->
      <t-button
        theme="primary"
        variant="outline"
        block
        @click="handleNewChat"
      >
        <template #icon><AddIcon /></template>
        {{ t('sidebar.newChat') }}
      </t-button>
    </div>

    <!-- Conversation List -->
    <div class="flex-1 overflow-y-auto py-2">
      <div
        v-if="chatStore.conversations.length === 0"
        class="px-4 py-8 text-center text-xs text-text-secondary"
      >
        {{ t('sidebar.emptyHistory') }}
      </div>

      <div
        v-for="conv in chatStore.conversations"
        :key="conv.id"
        class="conv-item group mx-2 mb-0.5 px-3 py-2.5 rounded-lg cursor-pointer flex items-start gap-2 transition-colors"
        :class="{
          'bg-brand/10 text-brand': currentConversationId === conv.id,
          'hover:bg-bg-hover': currentConversationId !== conv.id,
        }"
        @click="handleSelectConversation(conv.id)"
      >
        <div class="flex-1 min-w-0">
          <div
            class="text-sm font-medium truncate"
            :class="currentConversationId === conv.id ? 'text-brand' : 'text-text-primary'"
          >
            {{ conv.title || t('sidebar.newChat') }}
          </div>
          <div class="text-xs text-text-secondary mt-0.5">
            {{ formatDate(conv.created_at) }}
          </div>
        </div>
        <!-- Delete button (visible on hover) -->
        <t-button
          size="small"
          variant="text"
          shape="square"
          class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
          :class="{ '!opacity-100': currentConversationId === conv.id }"
          @click="(e: Event) => handleDeleteConversation(e, conv.id)"
        >
          <template #icon>
            <DeleteIcon class="text-text-secondary hover:text-error" />
          </template>
        </t-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.conv-sidebar {
  background: var(--td-bg-color-secondarycontainer, var(--td-bg-color-container));
}
</style>
