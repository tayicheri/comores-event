/// <reference types="vite/client" />
/// <reference types="vite-ssg" />

interface ImportMetaEnv {
  readonly VITE_API_URL?: string
  readonly VITE_APP_NAME?: string
  readonly VITE_SITE_URL?: string
  readonly VITE_DEV_PROXY_API?: string
  readonly VITE_STRIPE_PUBLISHABLE_KEY?: string
  readonly SSR?: boolean
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
