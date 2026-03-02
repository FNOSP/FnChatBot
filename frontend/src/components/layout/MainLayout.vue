<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { MoonIcon, SunnyIcon, TranslateIcon } from 'tdesign-icons-vue-next'
import Sidebar from './Sidebar.vue'
import ConversationSidebar from './ConversationSidebar.vue'
import { useTheme } from '../../composables/useTheme'

const route = useRoute()
const { t, locale } = useI18n()
const { toggleTheme, isDark } = useTheme()

// Show conversation sidebar only on chat routes (not settings)
const showConversationSidebar = computed(() => {
  return route.path === '/' || route.path.startsWith('/chat')
})

// Derive the page title from the current route
const pageTitle = computed(() => {
  const p = route.path
  if (p === '/' || p.startsWith('/chat')) return t('nav.chat')
  if (p.startsWith('/settings/general')) return t('nav.general')
  if (p.startsWith('/settings/sandbox')) return t('nav.sandbox')
  if (p.startsWith('/settings/models')) return t('nav.models')
  if (p.startsWith('/settings/mcp')) return t('nav.mcp')
  if (p.startsWith('/settings/skills')) return t('nav.skills')
  if (p.startsWith('/settings/users')) return t('nav.users')
  return t('settings.title')
})

// Language options
const langOptions = [
  { value: 'zh', label: '中文' },
  { value: 'en', label: 'English' },
  { value: 'ja', label: '日本語' },
]

const handleLangChange = (val: string) => {
  locale.value = val
  localStorage.setItem('locale', val)
}
</script>

<template>
  <div class="h-screen w-screen overflow-hidden flex bg-bg-primary">
    <!-- Panel 1: Navigation Sidebar -->
    <Sidebar />

    <!-- Panel 2: Conversation Sidebar (chat mode only) -->
    <ConversationSidebar v-if="showConversationSidebar" />

    <!-- Panel 3: Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Title Bar -->
      <header class="content-header h-12 shrink-0 flex items-center justify-between px-5 border-b border-border">
        <!-- Left: current page name -->
        <h1 class="text-sm font-semibold text-text-primary tracking-wide">
          {{ pageTitle }}
        </h1>

        <!-- Right: theme toggle + language selector -->
        <div class="flex items-center gap-2">
          <!-- Theme toggle -->
          <t-button
            variant="text"
            shape="square"
            size="small"
            @click="toggleTheme"
            :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
          >
            <template #icon>
              <MoonIcon v-if="!isDark" />
              <SunnyIcon v-else />
            </template>
          </t-button>

          <!-- Language selector -->
          <t-select
            :value="locale"
            :options="langOptions"
            size="small"
            style="width: 100px"
            @change="(val: string) => handleLangChange(val)"
          >
            <template #prefixIcon>
              <TranslateIcon />
            </template>
          </t-select>
        </div>
      </header>

      <!-- Content Slot -->
      <div class="flex-1 min-h-0 overflow-hidden bg-bg-card">
        <slot />
      </div>
    </div>
  </div>
</template>

<style scoped>
.content-header {
  background: var(--td-bg-color-container);
}
</style>
