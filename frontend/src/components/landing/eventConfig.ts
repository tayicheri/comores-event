export const EVENT_CONFIG = Object.freeze({
  name: 'Salon de la culture et de la gastronomie comorienne',
  year: 2027,
  dates: {
    vendredi: 'Vendredi 15 janvier',
    samedi: 'Samedi 16 janvier',
  } as const,
  dateShort: '15 & 16 janvier 2027',
  dateBadge: '15 & 16 JAN. 2027',
  dateKicker: '15 & 16 JANVIER 2027 À NANTERRE',
  dateOverlay: 'RENDEZ-VOUS LES 15 ET 16 JANVIER À NANTERRE',
  prices: { day: 20, pass: 30 } as const,
  currency: 'EUR',
  maxTickets: 20,
  venue: {
    name: 'Salle des congrès',
    address: '118 Rue du 8 Mai 1945, 92000 Nanterre',
    phones: ['07 77 16 86 84', '07 67 21 94 27'],
  },
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
