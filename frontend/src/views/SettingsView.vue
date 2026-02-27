<script setup lang="ts">
import { ref, computed, h, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Sidebar from '../components/layout/Sidebar.vue'
import GeneralSettings from '../components/settings/GeneralSettings.vue'
import SandboxSettings from '../components/settings/SandboxSettings.vue'
import ModelServices from '../components/settings/ModelServices.vue'
import MCPServers from '../components/settings/MCPServers.vue'
import SkillManagement from '../components/settings/SkillManagement.vue'
import UserManagement from '../components/settings/UserManagement.vue'
import { Sun, Moon } from 'lucide-vue-next'
import { 
  Layout as SemiLayout, 
  LayoutSider as SemiSider, 
  LayoutContent as SemiContent, 
  Nav as SemiNav 
} from '@kousum/semi-ui-vue'
import { 
  IconSetting, 
  IconShield,
  IconServer, 
  IconCloud, 
  IconCode 
} from '@kousum/semi-icons-vue'
import { useTheme } from '../composables/useTheme'

const currentTab = ref('general')
const isMobile = ref(false)
const { toggleTheme, isDark } = useTheme()

const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

const { t } = useI18n()

const menuItems = computed(() => [
  { itemKey: 'general', text: t('settings.general'), icon: IconSetting, component: GeneralSettings },
  { itemKey: 'sandbox', text: t('settings.sandbox'), icon: IconShield, component: SandboxSettings },
  { itemKey: 'models', text: t('settings.modelServices'), icon: IconServer, component: ModelServices },
  { itemKey: 'mcp', text: t('settings.mcpServers'), icon: IconCloud, component: MCPServers },
  { itemKey: 'skills', text: t('settings.skillManagement'), icon: IconCode, component: SkillManagement },
  { itemKey: 'users', text: t('settings.userManagement'), icon: IconCode, component: UserManagement }
])

const currentComponent = computed(() => {
  return menuItems.value.find(t => t.itemKey === currentTab.value)?.component
})

const onSelect = (data: { itemKey: string | number }) => {
  currentTab.value = String(data.itemKey)
}

const navItems = computed(() => menuItems.value.map(item => ({
  itemKey: item.itemKey,
  text: item.text,
  icon: h(item.icon)
})))
</script>

<template>
  <div class="flex h-screen w-full bg-background text-foreground">
    <!-- Main App Sidebar -->
    <Sidebar />
    
    <!-- Settings Layout -->
    <div class="flex-1 h-full overflow-hidden flex flex-col">
      <div class="h-12 border-b border-border flex items-center justify-end px-4">
        <button
          class="p-2 rounded-md hover:bg-secondary transition-colors"
          @click="toggleTheme"
          aria-label="Toggle theme"
        >
          <Moon v-if="isDark" class="w-4 h-4" />
          <Sun v-else class="w-4 h-4" />
        </button>
      </div>
      <SemiLayout class="flex-1" :style="{ flexDirection: isMobile ? 'column' : 'row' }">
        <SemiSider :style="{ width: isMobile ? '100%' : '240px' }">
          <SemiNav
            :selectedKeys="[currentTab]"
            :items="navItems"
            style="height: 100%"
            :mode="isMobile ? 'horizontal' : 'vertical'"
            @select="onSelect"
          >
            <template #header>
              <div class="p-4 pl-6" v-if="!isMobile">
                <h2 class="text-xl font-bold">{{ t('settings.title') }}</h2>
              </div>
            </template>
          </SemiNav>
        </SemiSider>
        <SemiContent class="h-full overflow-y-auto p-8">
          <component :is="currentComponent" />
        </SemiContent>
      </SemiLayout>
    </div>
  </div>
</template>
