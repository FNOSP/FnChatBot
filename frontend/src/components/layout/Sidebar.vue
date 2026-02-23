<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { MessageSquarePlus, Settings } from 'lucide-vue-next'

const router = useRouter()
const { t } = useI18n()

const navigateTo = (path: string) => {
  router.push(path)
}
</script>

<template>
  <aside class="w-64 bg-background text-foreground flex flex-col border-r border-border">
    <div class="p-4 border-b border-border">
      <h1 class="text-xl font-bold flex items-center gap-2">
        <span>FnChatBot</span>
      </h1>
    </div>

    <div class="p-2">
      <button 
        @click="navigateTo('/')"
        class="w-full flex items-center gap-2 p-2 rounded hover:bg-secondary transition-colors text-sm"
      >
        <MessageSquarePlus class="w-4 h-4" />
        {{ t('sidebar.newChat') }}
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-2">
      <div class="text-xs font-semibold text-muted-foreground mb-2 px-2 uppercase">{{ t('sidebar.recentChats') }}</div>
      <!-- History List Placeholder -->
      <div class="text-muted-foreground text-sm px-2 py-1 hover:bg-secondary rounded cursor-pointer">
        {{ t('sidebar.projectPlanning') }}
      </div>
    </div>

    <div class="p-2 border-t border-border">
      <button 
        @click="navigateTo('/settings')"
        class="w-full flex items-center gap-2 p-2 rounded hover:bg-secondary transition-colors text-sm"
      >
        <Settings class="w-4 h-4" />
        {{ t('settings.title') }}
      </button>
    </div>
  </aside>
</template>
