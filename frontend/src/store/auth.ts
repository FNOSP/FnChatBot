import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { http } from '../services/http'

export interface AuthUser {
  id: number
  username: string
  description?: string
  is_admin: boolean
  enabled: boolean
  must_change_password: boolean
}

interface LoginResponse {
  token: string
  user: AuthUser
  must_change_password: boolean
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('fnchatbot_token'))
  const currentUser = ref<AuthUser | null>(null)
  const mustChangePassword = ref(false)

  const isAuthenticated = computed(() => !!token.value && !!currentUser.value)
  const isAdmin = computed(() => !!currentUser.value?.is_admin)

  const setSession = (resp: LoginResponse) => {
    token.value = resp.token
    currentUser.value = resp.user
    mustChangePassword.value = resp.must_change_password
    localStorage.setItem('fnchatbot_token', resp.token)
  }

  const clearSession = () => {
    token.value = null
    currentUser.value = null
    mustChangePassword.value = false
    localStorage.removeItem('fnchatbot_token')
  }

  const login = async (username: string, password: string) => {
    const res = await http.post<LoginResponse>('/auth/login', { username, password })
    setSession(res.data)
    return res.data
  }

  const logout = () => {
    clearSession()
  }

  const fetchCurrentUser = async () => {
    if (!token.value) return
    try {
      const res = await http.get<AuthUser>('/auth/me')
      currentUser.value = res.data
      mustChangePassword.value = !!res.data.must_change_password
    } catch {
      clearSession()
    }
  }

  return {
    token,
    currentUser,
    mustChangePassword,
    isAuthenticated,
    isAdmin,
    login,
    logout,
    fetchCurrentUser
  }
})

