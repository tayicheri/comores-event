import { useHead } from '@unhead/vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { buildEventJsonLd } from '@/seo/eventJsonLd'

const appName = import.meta.env.VITE_APP_NAME || 'Salon de la culture'
const siteUrl = (import.meta.env.VITE_SITE_URL || '').replace(/\/$/, '')

/** Title + robots (noindex sur les routes meta.noindex). Canonical + JSON-LD Event sur la landing. */
export function usePageSeo() {
  const route = useRoute()

  const title = computed(() => {
    const metaTitle = route.meta.title
    if (!metaTitle) return appName
    if (route.meta.standaloneTitle) return metaTitle
    return `${metaTitle} · ${appName}`
  })

  useHead(() => {
    const head: {
      title: string
      meta: { name: string; content: string }[]
      script?: { key: string; type: string; children: string }[]
    } = {
      title: title.value,
      meta: [
        {
          name: 'robots',
          content: route.meta.noindex ? 'noindex, nofollow' : 'index, follow',
        },
      ],
    }

    // Canonical est dans index.html (__OG_URL__). JSON-LD uniquement sur la landing SSG.
    if (route.name === 'landing' && siteUrl) {
      head.script = [
        {
          key: 'ld-json-event',
          type: 'application/ld+json',
          children: JSON.stringify(buildEventJsonLd(siteUrl)),
        },
      ]
    }

    return head
  })
}
