export interface CurrentUser {
  id: string | number
  email: string
  display_name?: string | null
  role?: string
  roles?: string[]
  account_status?: string
  email_verified?: boolean
  created_at?: string
  [key: string]: unknown
}

export interface LoginCredentials {
  email: string
  password: string
}

export interface RegisterPayload {
  display_name?: string
  email: string
  password: string
}
