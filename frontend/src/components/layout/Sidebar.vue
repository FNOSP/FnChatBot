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
  if (route.path.startsWith('/chat')) return 'chat'
  if (route.path.startsWith('/settings')) return 'settings'
  return ''
})

const navigateTo = (path: string) => {
  router.push(path)
}

const handleNewChat = () => {
  router.push('/')
}
</script>

<template>
  <t-aside v-if="auth.isAuthenticated" width="240px">
    <t-menu :value="activeValue" theme="light" class="h-full">
      <template #logo>
        <div class="text-xl font-bold flex items-center gap-2 px-4 py-2">
          <span>FnChatBot</span>
        </div>
      </template>
      
      <t-menu-item value="new-chat" @click="handleNewChat">
        <template #icon>
          <AddIcon />
        </template>
        {{ t('sidebar.newChat') }}
      </t-menu-item>

      <t-menu-group :title="t('sidebar.recentChats')">
        <!-- History List Placeholder -->
        <t-menu-item value="chat-history-1" disabled>
          <template #icon>
            <ChatIcon />
          </template>
          {{ t('sidebar.projectPlanning') }}
        </t-menu-item>
      </t-menu-group>

      <template #operations>
        <t-menu-item value="settings" @click="navigateTo('/settings')">
          <template #icon>
            <SettingIcon />
          </template>
          {{ t('settings.title') }}
        </t-menu-item>
      </template>
    </t-menu>
  </t-aside>
</template>

