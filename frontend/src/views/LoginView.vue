<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'

const router = useRouter()
const auth = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const onSubmit = async () => {
  if (!username.value || !password.value) {
    error.value = 'Username and password are required.'
    return
  }
  error.value = ''
  loading.value = true
  try {
    const resp = await auth.login(username.value, password.value)
    if (resp.must_change_password) {
      router.push({ name: 'reset-password' })
    } else {
      router.push({ name: 'home' })
    }
  } catch (e: any) {
    error.value = e?.response?.data?.error || 'Login failed.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background text-foreground">
    <div class="w-full max-w-md bg-card border border-border rounded-lg p-8 shadow-md">
      <h1 class="text-2xl font-bold mb-6 text-center">FnChatBot</h1>
      <form @submit.prevent="onSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1">Username</label>
          <input
            v-model="username"
            type="text"
            class="w-full px-3 py-2 border border-border rounded-md bg-background"
            autocomplete="username"
          />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Password</label>
          <input
            v-model="password"
            type="password"
            class="w-full px-3 py-2 border border-border rounded-md bg-background"
            autocomplete="current-password"
          />
        </div>
        <p v-if="error" class="text-sm text-red-500">{{ error }}</p>
        <button
          type="submit"
          class="w-full py-2 mt-2 rounded-md bg-primary text-primary-foreground hover:opacity-90 disabled:opacity-60"
          :disabled="loading"
        >
          {{ loading ? 'Signing in...' : 'Sign In' }}
        </button>
      </form>
    </div>
  </div>
</template>

