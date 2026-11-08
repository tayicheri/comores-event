import { useHead } from '@unhead/vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const appName = import.meta.env.VITE_APP_NAME || 'Comores Event'

/** Title + robots (noindex sur les routes meta.noindex). */
export function usePageSeo() {
  const route = useRoute()

  const title = computed(() => {
    const metaTitle = route.meta.title
    if (!metaTitle) return appName
    if (route.meta.standaloneTitle) return metaTitle
    return `${metaTitle} · ${appName}`
  })

  useHead(() => ({
    title: title.value,
    meta: [
      {
        name: 'robots',
        content: route.meta.noindex ? 'noindex, nofollow' : 'index, follow',
      },
    ],
  }))
}
