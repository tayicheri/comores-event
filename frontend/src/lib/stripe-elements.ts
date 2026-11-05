import type { Stripe, StripeCardNumberElement, StripeElementStyle } from '@stripe/stripe-js'

export type CardNetwork = 'VISA' | 'MASTERCARD' | 'CARD'

export type CheckoutStripeFieldsApi = {
  getStripe: () => Stripe | null
  getCardNumber: () => StripeCardNumberElement | null
  fieldsComplete: () => boolean
}

/** Styles Elements face avant (fond forest). */
export const stripeCardElementStyleOnCard: StripeElementStyle = {
  base: {
    color: '#fffcf6',
    fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
    fontSize: '15px',
    fontWeight: '600',
    letterSpacing: '0.08em',
    '::placeholder': {
      color: 'rgba(255, 252, 246, 0.45)',
    },
    ...({ caretColor: 'transparent' } as object),
  },
  invalid: {
    color: '#f6e4b4',
  },
}

/** Styles CVC (bandeau blanc verso). */
export const stripeCardElementStyleOnCardCvc: StripeElementStyle = {
  base: {
    color: '#34314b',
    fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
    fontSize: '14px',
    fontWeight: '700',
    '::placeholder': {
      color: '#6e6a82',
    },
    ...({ caretColor: 'transparent' } as object),
  },
  invalid: {
    color: '#b91c1c',
  },
}

export function networkFromStripeBrand(brand: string): CardNetwork {
  if (brand === 'visa') return 'VISA'
  if (brand === 'mastercard') return 'MASTERCARD'
  return 'CARD'
}

export function stripeCardFieldsComplete(parts: {
  number: boolean
  expiry: boolean
  cvc: boolean
}): boolean {
  return parts.number && parts.expiry && parts.cvc
}
