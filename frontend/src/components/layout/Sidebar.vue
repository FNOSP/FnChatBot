<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ChatIcon, SettingIcon, AddIcon } from 'tdesign-icons-vue-next'
import { useAuthStore } from '../../store/auth'
import { computed } from 'vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const auth = useAuthStore()

const activeValue = computed(() => {
  // Map current route to left menu active key
  if (route.path.startsWith('/chat')) return 'chat'
  if (route.path.startsWith('/settings')) return 'settings'
  return ''
})

const navigateTo = (path: string) => {
  // Navigate to target path from sidebar actions
  router.push(path)
}

const handleNewChat = () => {
  // Start a fresh chat session from sidebar entry
  router.push('/')
}
</script>

<template>
  <t-aside
    v-if="auth.isAuthenticated"
    width="260px"
    class="bg-bg-secondary border-r border-border shadow-md flex flex-col"
  >
    <t-menu
      :value="activeValue"
      theme="dark"
      class="h-full !border-0 !bg-transparent flex-1 flex flex-col"
    >
      <template #logo>
        <div class="px-4 py-4 border-b border-border flex items-center justify-between">
          <div class="flex flex-col">
            <span class="text-lg font-semibold text-text-primary">FnChatBot</span>
            <span class="text-xs text-text-muted mt-0.5">
              {{ t('sidebar.appSubtitle') || 'AI Coding Assistant' }}
            </span>
          </div>
          <span class="px-2 py-0.5 text-[10px] rounded-full bg-brand/10 text-brand uppercase tracking-wide">
            Beta
          </span>
        </div>
      </template>
      
      <t-menu-item
        value="new-chat"
        @click="handleNewChat"
        class="mx-3 mt-4 mb-1 rounded-lg bg-brand text-text-inverse hover:bg-brand-light transition-colors"
      >
        <template #icon>
          <AddIcon />
        </template>
        {{ t('sidebar.newChat') }}
      </t-menu-item>

      <div class="flex-1 min-h-0 flex flex-col mt-2">
        <t-menu-group :title="t('sidebar.recentChats')" class="flex-1 min-h-0">
          <!-- History List Placeholder -->
          <div class="px-3 pb-3 text-xs text-text-muted">
            {{ t('sidebar.emptyHistory') || 'Your recent chats will appear here.' }}
          </div>
          <t-menu-item value="chat-history-1" disabled class="opacity-60">
            <template #icon>
              <ChatIcon />
            </template>
            {{ t('sidebar.projectPlanning') }}
          </t-menu-item>
        </t-menu-group>
      </div>

      <template #operations>
        <t-menu-item
          value="settings"
          @click="navigateTo('/settings')"
          class="border-t border-border mt-2 pt-3"
        >
          <template #icon>
            <SettingIcon />
          </template>
          {{ t('settings.title') }}
        </t-menu-item>
      </template>
    </t-menu>
  </t-aside>
</template>

