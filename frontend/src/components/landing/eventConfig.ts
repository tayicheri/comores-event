export const EVENT_CONFIG = Object.freeze({
  name: 'Salon de la culture et de la gastronomie comorienne',
  year: 2027,
  dates: {
    samedi: 'Samedi 16 janvier',
    dimanche: 'Dimanche 17 janvier',
  } as const,
  prices: { day: 20, pass: 30 } as const,
  currency: 'EUR',
  maxTickets: 20,
})

export type TicketFormula = keyof typeof EVENT_CONFIG.prices
export type TicketDay = keyof typeof EVENT_CONFIG.dates

export function formatEuro(value: number): string {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: EVENT_CONFIG.currency,
    maximumFractionDigits: 0,
  }).format(value)
}
