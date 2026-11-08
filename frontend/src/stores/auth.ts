import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import * as authApi from '@/api/auth'
import { registerAuthHooks } from '@/api/client'
import type { CurrentUser, LoginCredentials } from '@/types/auth'

function userIsAdmin(user: CurrentUser | null): boolean {
  if (!user) return false
  if (user.role === 'admin') return true
  return Array.isArray(user.roles) && user.roles.includes('admin')
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<CurrentUser | null>(null)
  const loading = ref(false)
  const isInitialized = ref(false)

  const isAuthenticated = computed(() => user.value !== null)
  const isAdmin = computed(() => userIsAdmin(user.value))
  const displayName = computed(
    () => (typeof user.value?.display_name === 'string' ? user.value.display_name.trim() : '') || user.value?.email || '',
  )

  function clearLocalSession() {
    user.value = null
  }

  async function fetchMe(): Promise<CurrentUser | null> {
    try {
      user.value = await authApi.me()
    } finally {
      isInitialized.value = true
    }
    return user.value
  }

  async function init() {
    if (isInitialized.value) return
    try {
      await fetchMe()
    } catch {
      user.value = null
    }
  }

  async function login(credentials: LoginCredentials) {
    loading.value = true
    try {
      await authApi.login(credentials)
      const nextUser = await authApi.me()
      if (!nextUser) {
        throw new Error('Session créée mais utilisateur introuvable.')
      }
      user.value = nextUser
      isInitialized.value = true
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // la session locale est purgée quoi qu'il arrive
    } finally {
      clearLocalSession()
    }
  }

  registerAuthHooks({
    isAuthenticated: () => user.value !== null,
    onSessionExpired: () => {
      clearLocalSession()
      void import('@/router').then(({ getRouter }) => {
        const router = getRouter()
        if (!router) {
          if (typeof window !== 'undefined') {
            window.location.assign('/lekeza/login')
          }
          return
        }
        if (router.currentRoute.value.name !== 'lekeza-login') {
          void router.replace({
            name: 'lekeza-login',
            query: { redirect: router.currentRoute.value.fullPath },
          })
        }
      })
    },
  })

  return {
    user,
    loading,
    isInitialized,
    isAuthenticated,
    isAdmin,
    displayName,
    init,
    fetchMe,
    login,
    logout,
    clearLocalSession,
  }
})
