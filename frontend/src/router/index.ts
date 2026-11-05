import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

declare module 'vue-router' {
  interface RouteMeta {
    public?: boolean
    requiresAdmin?: boolean
    title?: string
  }
}

const appName = import.meta.env.VITE_APP_NAME || 'Comores Event'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: () => import('@/views/LandingView.vue'),
      meta: { public: true, title: 'L’esprit des îles' },
    },
    {
      path: '/lekeza/login',
      name: 'lekeza-login',
      component: () => import('@/views/LoginView.vue'),
      meta: { public: true, title: 'Connexion admin' },
    },
    {
      path: '/lekeza',
      name: 'lekeza',
      component: () => import('@/views/AdminComingSoonView.vue'),
      meta: { requiresAdmin: true, title: 'Admin' },
    },
    { path: '/login', redirect: { name: 'lekeza-login' } },
    { path: '/:pathMatch(.*)*', redirect: { name: 'landing' } },
  ],
  scrollBehavior(to) {
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    return { top: 0 }
  },
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  await auth.init()

  if (to.name === 'lekeza-login') {
    if (auth.isAdmin) return { name: 'lekeza' }
    if (auth.isAuthenticated) return { name: 'landing' }
    return true
  }

  if (to.meta.public) {
    return true
  }

  if (to.meta.requiresAdmin) {
    if (!auth.isAuthenticated) {
      return {
        name: 'lekeza-login',
        query: to.fullPath !== '/lekeza' ? { redirect: to.fullPath } : undefined,
      }
    }
    if (!auth.isAdmin) {
      return { name: 'landing' }
    }
    return true
  }

  return { name: 'landing' }
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} · ${appName}` : appName
})

export default router
