import { ViteSSG } from 'vite-ssg'
import { createPinia } from 'pinia'
import App from '@/App.vue'
import { routes, scrollBehavior, setRouter, setupRouterGuards } from '@/router'
import '@/style.css'

export const createApp = ViteSSG(
  App,
  { routes, scrollBehavior },
  ({ app, router, isClient }) => {
    app.use(createPinia())
    setRouter(router)
    if (isClient) {
      setupRouterGuards(router)
    }
  },
)
