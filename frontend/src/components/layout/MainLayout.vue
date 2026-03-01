<script setup lang="ts">
import { MoonIcon, SunnyIcon } from 'tdesign-icons-vue-next'
import { useI18n } from 'vue-i18n'
import Sidebar from './Sidebar.vue'
import SettingsView from '../../views/SettingsView.vue'
import { useTheme } from '../../composables/useTheme'
import { useSettingsDrawer } from '../../composables/useSettingsDrawer'

const { t } = useI18n()
const { toggleTheme, isDark } = useTheme()
const { settingsVisible, closeSettings } = useSettingsDrawer()
</script>

<template>
  <t-layout class="h-screen w-screen overflow-hidden bg-bg-primary">
    <Sidebar />
    <t-layout class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <header class="h-14 shrink-0 flex items-center justify-between px-6 bg-bg-card">
        <div></div>
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

    <t-drawer
      v-model:visible="settingsVisible"
      :header="t('settings.title')"
      placement="right"
      size="75%"
      :footer="false"
      :close-on-overlay-click="true"
      :show-overlay="true"
      @close="closeSettings"
    >
      <SettingsView />
    </t-drawer>
  </t-layout>
</template>
