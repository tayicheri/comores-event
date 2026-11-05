import axios, { type AxiosError, type AxiosRequestConfig, type InternalAxiosRequestConfig } from 'axios'

declare module 'axios' {
  interface AxiosRequestConfig {
    csrfRetried?: boolean
    refreshRetried?: boolean
  }
}

const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const MUTATING_METHODS = ['POST', 'PUT', 'PATCH', 'DELETE']

export const CSRF_COOKIE = 'XSRF-TOKEN'
export const CSRF_HEADER = 'X-XSRF-TOKEN'

export const apiClient = axios.create({
  baseURL: apiBaseUrl,
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
})

export function readCookie(name: string): string | null {
  const match = document.cookie.split('; ').find((cookie) => cookie.startsWith(`${name}=`))
  return match ? decodeURIComponent(match.split('=').slice(1).join('=')) : null
}

function clearClientCookie(name: string) {
  document.cookie = `${name}=; Max-Age=0; path=/`
  const hostname = window.location.hostname
  if (hostname && hostname !== 'localhost') {
    document.cookie = `${name}=; Max-Age=0; path=/; domain=${hostname}`
    document.cookie = `${name}=; Max-Age=0; path=/; domain=.${hostname}`
  }
}

/** Callbacks injectés par le store auth pour éviter une dépendance circulaire. */
interface AuthHooks {
  isAuthenticated: () => boolean
  onSessionExpired: () => void
}

let authHooks: AuthHooks | null = null

export function registerAuthHooks(hooks: AuthHooks) {
  authHooks = hooks
}

/** Single-flight : plusieurs 401 concurrents partagent un seul POST /api/auth/refresh. */
let refreshInFlight: Promise<boolean> | null = null

function tryRefreshSession(): Promise<boolean> {
  if (refreshInFlight) {
    return refreshInFlight
  }
  refreshInFlight = (async () => {
    try {
      await apiClient.post('/api/auth/refresh', null, { refreshRetried: true } satisfies AxiosRequestConfig)
      return true
    } catch {
      return false
    } finally {
      refreshInFlight = null
    }
  })()
  return refreshInFlight
}

function isAuthFlowUrl(url: string): boolean {
  return (
    url.includes('/api/auth/login') ||
    url.includes('/api/auth/register') ||
    url.includes('/api/auth/logout') ||
    url.includes('/api/auth/refresh')
  )
}

function isCsrfFailure(status: number, data: unknown): boolean {
  if (status !== 403) {
    return false
  }
  const raw = typeof data === 'string' ? data : JSON.stringify(data ?? '')
  return /csrf|xsrf/i.test(raw)
}

apiClient.interceptors.request.use((config) => {
  const method = (config.method ?? 'get').toUpperCase()
  if (MUTATING_METHODS.includes(method)) {
    const token = readCookie(CSRF_COOKIE)
    if (token) {
      config.headers.set(CSRF_HEADER, token)
    }
  }
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<{ message?: string; error?: string }>) => {
    let message = 'Une erreur inattendue est survenue.'
    const failedConfig = error.config as InternalAxiosRequestConfig | undefined
    const requestUrl = failedConfig?.url ?? ''

    if (error.response) {
      const { status, data } = error.response
      message = data?.message ?? data?.error ?? message

      if (status === 401) {
        const isMeProbe = requestUrl.includes('/api/auth/me')
        const canTryRefresh = Boolean(failedConfig) && !failedConfig?.refreshRetried && !isAuthFlowUrl(requestUrl)

        if (canTryRefresh && failedConfig) {
          if (await tryRefreshSession()) {
            failedConfig.refreshRetried = true
            return apiClient.request(failedConfig)
          }
        }

        if (!isAuthFlowUrl(requestUrl)) {
          if (authHooks?.isAuthenticated()) {
            authHooks.onSessionExpired()
          } else if (isMeProbe) {
            clearClientCookie(CSRF_COOKIE)
          }
        }
      }

      if (isCsrfFailure(status, data)) {
        message = 'Session de sécurité expirée. Veuillez réessayer.'
        if (failedConfig && !failedConfig.csrfRetried && !requestUrl.includes('/health')) {
          clearClientCookie(CSRF_COOKIE)
          try {
            await apiClient.get('/health')
            failedConfig.csrfRetried = true
            return apiClient.request(failedConfig)
          } catch {
            // on retombe sur le message CSRF ci-dessus
          }
        }
      }
    } else if (error.request) {
      message = 'Impossible de contacter le serveur. Vérifiez votre connexion.'
    }

    error.message = message
    return Promise.reject(error)
  },
)

export function getErrorMessage(error: unknown, fallback = 'Une erreur inattendue est survenue.'): string {
  if (error instanceof Error && error.message) {
    return error.message
  }
  return fallback
}
