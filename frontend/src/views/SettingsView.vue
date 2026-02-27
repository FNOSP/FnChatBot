<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import GeneralSettings from '../components/settings/GeneralSettings.vue'
import SandboxSettings from '../components/settings/SandboxSettings.vue'
import ModelServices from '../components/settings/ModelServices.vue'
import MCPServers from '../components/settings/MCPServers.vue'
import SkillManagement from '../components/settings/SkillManagement.vue'
import UserManagement from '../components/settings/UserManagement.vue'
import { SettingIcon, SecuredIcon, ServerIcon, CloudIcon, CodeIcon, UserIcon } from 'tdesign-icons-vue-next'

const currentTab = ref('general')
const { t } = useI18n()

const menuItems = computed(() => [
  { value: 'general', label: t('settings.general'), icon: SettingIcon, component: GeneralSettings },
  { value: 'sandbox', label: t('settings.sandbox'), icon: SecuredIcon, component: SandboxSettings },
  { value: 'models', label: t('settings.modelServices'), icon: ServerIcon, component: ModelServices },
  { value: 'mcp', label: t('settings.mcpServers'), icon: CloudIcon, component: MCPServers },
  { value: 'skills', label: t('settings.skillManagement'), icon: CodeIcon, component: SkillManagement },
  { value: 'users', label: t('settings.userManagement'), icon: UserIcon, component: UserManagement }
])
</script>

<template>
  <div class="h-full flex flex-col p-4 bg-bg-primary">
    <div class="mb-4">
      <h2 class="text-2xl font-bold">{{ t('settings.title') }}</h2>
    </div>
    
    <t-tabs v-model="currentTab" placement="left" theme="normal" class="flex-1 bg-bg-card rounded-lg shadow-sm border border-border">
      <t-tab-panel 
        v-for="item in menuItems" 
        :key="item.value" 
        :value="item.value"
        :label="item.label"
      >
        <template #label>
          <div class="flex items-center gap-2">
            <component :is="item.icon" />
            <span>{{ item.label }}</span>
          </div>
        </template>
        <div class="p-6 h-full overflow-y-auto">
          <component :is="item.component" />
        </div>
      </t-tab-panel>
    </t-tabs>
  </div>
</template>

