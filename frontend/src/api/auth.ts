import { apiClient } from '@/api/client'
import type { CurrentUser, LoginCredentials, RegisterPayload } from '@/types/auth'

type UserEnvelope = CurrentUser | { user: CurrentUser } | { data: CurrentUser }

function unwrapUser(payload: UserEnvelope): CurrentUser {
  const envelope = payload as { user?: CurrentUser; data?: CurrentUser }
  return envelope.user ?? envelope.data ?? (payload as CurrentUser)
}

export async function login(credentials: LoginCredentials): Promise<void> {
  await apiClient.post('/api/auth/login', credentials)
}

export async function register(payload: RegisterPayload): Promise<void> {
  await apiClient.post('/api/auth/register', payload)
}

export async function logout(): Promise<void> {
  await apiClient.post('/api/auth/logout')
}

export async function refresh(): Promise<void> {
  await apiClient.post('/api/auth/refresh')
}

/** Retourne l'utilisateur courant, ou null si aucune session (401). */
export async function me(): Promise<CurrentUser | null> {
  try {
    const { data } = await apiClient.get<UserEnvelope>('/api/auth/me')
    return unwrapUser(data)
  } catch (error) {
    const status = (error as { response?: { status?: number } }).response?.status
    if (status === 401) {
      return null
    }
    throw error
  }
}
