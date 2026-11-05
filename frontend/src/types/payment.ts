import type { TicketDay, TicketFormula } from '@/components/landing/eventConfig'

export type CheckoutSelection = {
  formula: TicketFormula
  day: TicketDay
  quantity: number
  unit: number
  total: number
}

export type BuyerForm = {
  first_name: string
  last_name: string
  email: string
  city: string
}

export type PaymentIntentResponse = {
  payment_intent_id: string
  client_secret: string
  status: string
  amount_cents: number
  currency: string
}

export type BilletStatus = 'valid' | 'used'

export type Billet = {
  id: string
  purchased_at: string
  first_name: string
  last_name: string
  email: string
  city: string
  quantity: number
  day_vendredi: boolean
  day_samedi: boolean
  comments: string | null
  status: BilletStatus
  amount_cents: number
  stripe_payment_intent_id: string
}

export type BilletStats = {
  total_count: number
  valid_count: number
  used_count: number
  revenue_cents: number
  vendredi_places: number
  samedi_places: number
}

export type BilletListResponse = {
  billets: Billet[]
  stats: BilletStats
}
