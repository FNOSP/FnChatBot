<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { http } from '../services/http'
import { useAuthStore } from '../store/auth'

const router = useRouter()
const auth = useAuthStore()

// Hide old password field when forced to change initial password.
const requireOldPassword = computed(() => !auth.mustChangePassword)

const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

const onSubmit = async () => {
  error.value = ''
  success.value = ''

  if (!newPassword.value || !confirmPassword.value) {
    error.value = 'New password and confirmation are required.'
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    error.value = 'Passwords do not match.'
    return
  }

  loading.value = true
  try {
    await http.post('/auth/reset-password', {
      old_password: requireOldPassword.value ? oldPassword.value : undefined,
      new_password: newPassword.value,
      new_password_confirm: confirmPassword.value
    })
    auth.mustChangePassword = false
    success.value = 'Password updated successfully.'
    router.push({ name: 'home' })
  } catch (e: any) {
    error.value = e?.response?.data?.error || 'Failed to reset password.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background text-foreground">
    <div class="w-full max-w-md bg-card border border-border rounded-lg p-8 shadow-md">
      <h1 class="text-2xl font-bold mb-4 text-center">Reset Password</h1>
      <p class="text-sm text-muted-foreground mb-4 text-center">
        You must change the initial password before using the system.
      </p>
      <form @submit.prevent="onSubmit" class="space-y-4">
        <div v-if="requireOldPassword">
          <label class="block text-sm font-medium mb-1">Old Password</label>
          <input
            v-model="oldPassword"
            type="password"
            class="w-full px-3 py-2 border border-border rounded-md bg-background"
            autocomplete="current-password"
          />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">New Password</label>
          <input
            v-model="newPassword"
            type="password"
            class="w-full px-3 py-2 border border-border rounded-md bg-background"
            autocomplete="new-password"
          />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Confirm New Password</label>
          <input
            v-model="confirmPassword"
            type="password"
            class="w-full px-3 py-2 border border-border rounded-md bg-background"
            autocomplete="new-password"
          />
        </div>

        <p v-if="error" class="text-sm text-red-500">{{ error }}</p>
        <p v-if="success" class="text-sm text-green-500">{{ success }}</p>

        <button
          type="submit"
          class="w-full py-2 mt-2 rounded-md bg-primary text-primary-foreground hover:opacity-90 disabled:opacity-60"
          :disabled="loading"
        >
          {{ loading ? 'Updating...' : 'Update Password' }}
        </button>
      </form>
    </div>
  </div>
</template>

