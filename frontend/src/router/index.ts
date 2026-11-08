import type { RouteRecordRaw, Router, RouterScrollBehavior } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

let routerInstance: Router | null = null

export function setRouter(router: Router) {
  routerInstance = router
}

export function getRouter(): Router | null {
  return routerInstance
}

declare module 'vue-router' {
  interface RouteMeta {
    public?: boolean
    requiresAdmin?: boolean
    title?: string
    standaloneTitle?: boolean
    /** Si true, meta robots noindex,nofollow (admin). */
    noindex?: boolean
  }
}

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'landing',
    component: () => import('@/views/LandingView.vue'),
    meta: {
      public: true,
      title: 'Salon de la culture et de la gastronomie comorienne en France',
      standaloneTitle: true,
    },
  },
  {
    path: '/lekeza/login',
    name: 'lekeza-login',
    component: () => import('@/views/LoginView.vue'),
    meta: { public: true, title: 'Connexion admin', noindex: true },
  },
  {
    path: '/lekeza',
    name: 'lekeza',
    component: () => import('@/views/AdminDashboardView.vue'),
    meta: { requiresAdmin: true, title: 'Admin', noindex: true },
  },
  { path: '/login', redirect: { name: 'lekeza-login' } },
  { path: '/:pathMatch(.*)*', redirect: { name: 'landing' } },
]

export const scrollBehavior: RouterScrollBehavior = (to) => {
  if (to.hash) {
    return { el: to.hash, behavior: 'smooth' }
  }
  return { top: 0 }
}

/** Guards auth — client uniquement (pas pendant le prérendu SSG). */
export function setupRouterGuards(router: Router) {
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
}
