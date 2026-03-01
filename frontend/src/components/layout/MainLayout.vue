<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { MoonIcon, SunnyIcon } from 'tdesign-icons-vue-next'
import { useI18n } from 'vue-i18n'
import Sidebar from './Sidebar.vue'
import { useTheme } from '../../composables/useTheme'

const route = useRoute()
const { t } = useI18n()
const { toggleTheme, isDark } = useTheme()

const pageTitle = computed(() => {
  if (route.path === '/settings') return t('settings.title')
  return 'FnChatBot'
})

const pageSubtitle = computed(() => {
  if (route.path === '/settings') return ''
  return t('chat.conversation') || 'Conversation'
})
</script>

<template>
  <t-layout class="h-screen w-screen overflow-hidden bg-bg-primary">
    <Sidebar />
    <t-layout class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <header class="h-14 shrink-0 border-b border-border flex items-center justify-between px-6 bg-bg-card">
        <div class="flex items-center gap-3">
          <div class="flex flex-col">
            <span v-if="pageSubtitle" class="text-xs text-text-secondary uppercase tracking-wide leading-tight">
              {{ pageSubtitle }}
            </span>
            <span class="text-base font-semibold text-text-primary leading-tight">
              {{ pageTitle }}
            </span>
          </div>
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
      <div class="flex-1 flex flex-col min-w-0 overflow-hidden bg-bg-card">
        <slot />
      </div>
    </t-layout>
  </t-layout>
</template>
