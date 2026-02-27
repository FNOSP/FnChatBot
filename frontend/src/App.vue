<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './store/auth'
import MainLayout from './components/layout/MainLayout.vue'

const router = useRouter()
const auth = useAuthStore()
// Show loading until router has finished initial navigation (including auth redirect)
const isRouterReady = ref(false)

onMounted(async () => {
  await router.isReady()
  isRouterReady.value = true
})
</script>

<template>
  <div v-if="!isRouterReady" class="min-h-screen flex items-center justify-center bg-bg-primary text-text-primary">
    <div class="text-center">
      <div class="inline-block w-8 h-8 border-2 border-brand border-t-transparent rounded-full animate-spin mb-2" />
      <p>Loading...</p>
    </div>
  </div>
  <template v-else>
    <MainLayout v-if="auth.isAuthenticated">
      <router-view />
    </MainLayout>
    <router-view v-else />
  </template>
</template>

<style>
/* Global styles if needed */
</style>
