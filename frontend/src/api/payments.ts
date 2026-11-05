import { apiClient } from '@/api/client'
import type { BuyerForm, CheckoutSelection, PaymentIntentResponse, Billet } from '@/types/payment'

export async function createPaymentIntent(
  buyer: BuyerForm,
  selection: CheckoutSelection,
): Promise<PaymentIntentResponse> {
  const { data } = await apiClient.post<PaymentIntentResponse>('/api/payments/intent', {
    first_name: buyer.first_name.trim(),
    last_name: buyer.last_name.trim(),
    email: buyer.email.trim(),
    city: buyer.city.trim(),
    formula: selection.formula,
    day: selection.formula === 'day' ? selection.day : undefined,
    quantity: selection.quantity,
  })
  return data
}

export async function confirmPayment(paymentIntentId: string): Promise<Billet> {
  const { data } = await apiClient.post<{ billet: Billet }>('/api/payments/confirm', {
    payment_intent_id: paymentIntentId,
  })
  return data.billet
}

/** Pose le cookie CSRF anonyme (GET sûr côté API). */
export async function ensureCsrfCookie(): Promise<void> {
  await apiClient.get('/health')
}
