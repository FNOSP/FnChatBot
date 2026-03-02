import { createRouter, createWebHistory } from 'vue-router'
import ChatView from '../views/ChatView.vue'
import SettingsLayout from '../components/layout/SettingsLayout.vue'
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
      component: SettingsLayout,
      redirect: { name: 'settings-general' },
      children: [
        {
          path: 'general',
          name: 'settings-general',
          component: () => import('../components/settings/GeneralSettings.vue')
        },
        {
          path: 'sandbox',
          name: 'settings-sandbox',
          component: () => import('../components/settings/SandboxSettings.vue')
        },
        {
          path: 'models',
          name: 'settings-models',
          component: () => import('../components/settings/ModelServices.vue')
        },
        {
          path: 'mcp',
          name: 'settings-mcp',
          component: () => import('../components/settings/MCPServers.vue')
        },
        {
          path: 'skills',
          name: 'settings-skills',
          component: () => import('../components/settings/SkillManagement.vue')
        },
        {
          path: 'users',
          name: 'settings-users',
          component: () => import('../components/settings/UserManagement.vue')
        },
      ]
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
