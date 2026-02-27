import { createRouter, createWebHistory } from 'vue-router'
import ChatView from '../views/ChatView.vue'
import { useAuthStore } from '../store/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ChatView
    },
    {
      path: '/chat/:id',
      name: 'chat',
      component: ChatView
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/SettingsView.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: () => import('../views/ResetPasswordView.vue')
    }
  ]
})

router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore()

  // Public routes.
  if (to.name === 'login') {
    if (auth.isAuthenticated && !auth.mustChangePassword) {
      return next({ name: 'home' })
    }
    return next()
  }

  // Require authentication for all other routes.
  if (!auth.isAuthenticated) {
    await auth.fetchCurrentUser()
  }

  if (!auth.isAuthenticated) {
    return next({ name: 'login' })
  }

  // Force password reset if required.
  if (auth.mustChangePassword && to.name !== 'reset-password') {
    return next({ name: 'reset-password' })
  }

  next()
})

export default router
