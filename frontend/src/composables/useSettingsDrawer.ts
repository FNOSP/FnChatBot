import { ref } from 'vue'

const settingsVisible = ref(false)

export function useSettingsDrawer() {
  const toggleSettings = () => {
    settingsVisible.value = !settingsVisible.value
  }

  const openSettings = () => {
    settingsVisible.value = true
  }

  const closeSettings = () => {
    settingsVisible.value = false
  }

  return {
    settingsVisible,
    toggleSettings,
    openSettings,
    closeSettings,
  }
}
