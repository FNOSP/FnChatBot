import { ref, watchEffect, onMounted, onUnmounted, computed } from 'vue'

type Theme = 'light' | 'dark'
type ThemePreference = Theme | 'system'

const theme = ref<Theme>('light')
const preference = ref<ThemePreference>('system')
let mediaQuery: MediaQueryList | null = null
let initialized = false

const getStoredPreference = (): ThemePreference => {
  const saved = localStorage.getItem('theme')
  if (saved === 'light' || saved === 'dark') return saved
  return 'system'
}

const resolveTheme = () => {
  const resolved = preference.value === 'system'
    ? (mediaQuery?.matches ? 'dark' : 'light')
    : preference.value
  theme.value = resolved
}

const applyTheme = (t: Theme) => {
  document.documentElement.classList.remove('light', 'dark')
  document.documentElement.classList.add(t)
  document.documentElement.setAttribute('theme-mode', t)
  document.documentElement.setAttribute('data-theme', t)
  document.body.setAttribute('theme-mode', t)
  document.body.setAttribute('data-theme', t)
}

const syncStorage = () => {
  if (preference.value === 'system') {
    localStorage.removeItem('theme')
    return
  }
  localStorage.setItem('theme', preference.value)
}

const setupMediaListener = () => {
  if (!mediaQuery) return
  const handler = () => {
    if (preference.value === 'system') {
      resolveTheme()
      applyTheme(theme.value)
    }
  }
  mediaQuery.addEventListener('change', handler)
  return () => mediaQuery?.removeEventListener('change', handler)
}

const initTheme = () => {
  if (initialized) return
  initialized = true
  mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  preference.value = getStoredPreference()
  resolveTheme()
  applyTheme(theme.value)
  syncStorage()
}

export function useTheme() {
  let cleanup: (() => void) | undefined

  onMounted(() => {
    initTheme()
    cleanup = setupMediaListener()
  })

  onUnmounted(() => {
    cleanup?.()
  })

  const toggleTheme = () => {
    preference.value = theme.value === 'dark' ? 'light' : 'dark'
    resolveTheme()
    applyTheme(theme.value)
    syncStorage()
  }

  watchEffect(() => {
    if (!initialized) return
    applyTheme(theme.value)
  })

  return {
    theme,
    toggleTheme,
    isDark: computed(() => theme.value === 'dark'),
  }
}
