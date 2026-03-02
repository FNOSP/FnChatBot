<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  ChatIcon,
  RobotIcon,
  SettingIcon,
  SecuredIcon,
  ServerIcon,
  CloudIcon,
  CodeIcon,
  UserIcon,
  AppIcon,
} from 'tdesign-icons-vue-next'
import { useAuthStore } from '../../store/auth'
import { useTheme } from '../../composables/useTheme'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const { isAdmin } = useAuthStore()
const { isDark } = useTheme()

const menuTheme = computed(() => (isDark.value ? 'dark' : 'light'))

// Compute the active menu value from current route
const activeValue = computed(() => {
  const p = route.path
  if (p.startsWith('/settings/general')) return 'settings-general'
  if (p.startsWith('/settings/sandbox')) return 'settings-sandbox'
  if (p.startsWith('/settings/models')) return 'settings-models'
  if (p.startsWith('/settings/mcp')) return 'settings-mcp'
  if (p.startsWith('/settings/skills')) return 'settings-skills'
  if (p.startsWith('/settings/users')) return 'settings-users'
  if (p.startsWith('/chat') || p === '/') return 'chat'
  return ''
})

// Which submenus should be expanded by default
const defaultExpanded = computed(() => {
  const p = route.path
  const expanded: string[] = ['section-chat']
  if (p.startsWith('/settings/skills') || p.startsWith('/settings/mcp')) {
    expanded.push('section-agent')
  }
  if (
    p.startsWith('/settings/general') ||
    p.startsWith('/settings/sandbox') ||
    p.startsWith('/settings/models') ||
    p.startsWith('/settings/users')
  ) {
    expanded.push('section-settings')
  }
  return expanded
})

const navigate = (path: string) => {
  router.push(path)
}
</script>

<template>
  <div
    class="nav-sidebar flex flex-col border-r border-border shrink-0"
    style="width: 200px;"
  >
    <!-- Logo -->
    <div class="px-4 py-4 flex items-center gap-2 shrink-0 border-b border-border">
      <AppIcon class="text-brand text-xl" />
      <span class="text-base font-bold text-text-primary tracking-tight">FnChatBot</span>
      <span class="ml-auto px-1.5 py-0.5 text-[9px] rounded-full bg-brand/10 text-brand uppercase tracking-wide shrink-0">
        Beta
      </span>
    </div>

    <!-- Navigation Menu -->
    <t-menu
      :value="activeValue"
      :theme="menuTheme"
      :default-expanded="defaultExpanded"
      collapsed-width="0"
      class="nav-menu flex-1 !border-0 !bg-transparent overflow-y-auto"
    >
      <!-- Chat Section -->
      <t-submenu value="section-chat">
        <template #icon><ChatIcon /></template>
        <template #title>{{ t('sidebar.sectionChat') }}</template>
        <t-menu-item value="chat" @click="navigate('/')">
          <template #icon><ChatIcon /></template>
          {{ t('sidebar.chat') }}
        </t-menu-item>
      </t-submenu>

      <!-- Agent Section -->
      <t-submenu value="section-agent">
        <template #icon><RobotIcon /></template>
        <template #title>{{ t('sidebar.sectionAgent') }}</template>
        <t-menu-item value="settings-skills" @click="navigate('/settings/skills')">
          <template #icon><CodeIcon /></template>
          {{ t('settings.skillManagement') }}
        </t-menu-item>
        <t-menu-item value="settings-mcp" @click="navigate('/settings/mcp')">
          <template #icon><CloudIcon /></template>
          {{ t('settings.mcpServers') }}
        </t-menu-item>
      </t-submenu>

      <!-- Settings Section -->
      <t-submenu value="section-settings">
        <template #icon><SettingIcon /></template>
        <template #title>{{ t('sidebar.sectionSettings') }}</template>
        <t-menu-item value="settings-general" @click="navigate('/settings/general')">
          <template #icon><SettingIcon /></template>
          {{ t('settings.general') }}
        </t-menu-item>
        <t-menu-item value="settings-sandbox" @click="navigate('/settings/sandbox')">
          <template #icon><SecuredIcon /></template>
          {{ t('settings.sandbox') }}
        </t-menu-item>
        <t-menu-item value="settings-models" @click="navigate('/settings/models')">
          <template #icon><ServerIcon /></template>
          {{ t('settings.modelServices') }}
        </t-menu-item>
        <t-menu-item
          v-if="isAdmin"
          value="settings-users"
          @click="navigate('/settings/users')"
        >
          <template #icon><UserIcon /></template>
          {{ t('settings.userManagement') }}
        </t-menu-item>
      </t-submenu>
    </t-menu>
  </div>
</template>

<style scoped>
.nav-sidebar {
  background: var(--td-bg-color-container);
}

:deep(.nav-menu) {
  width: 100% !important;
}

:deep(.t-menu__item) {
  border-radius: 6px;
  margin: 1px 6px;
  width: calc(100% - 12px) !important;
}

:deep(.t-submenu__title) {
  border-radius: 6px;
  margin: 1px 6px;
  width: calc(100% - 12px) !important;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.02em;
  opacity: 0.6;
  text-transform: uppercase;
}

:deep(.t-menu__item--active) {
  background: var(--td-brand-color-light) !important;
  color: var(--td-brand-color) !important;
}

:deep(.t-submenu .t-menu__item) {
  padding-left: 36px !important;
  font-size: 13px;
}
</style>
