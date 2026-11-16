export const EVENT_CONFIG = Object.freeze({
  name: 'Salon de la culture et de la gastronomie comorienne',
  year: 2027,
  dates: {
    // Clés API backend inchangées (day_vendredi / day_samedi) → 1er / 2e jour.
    vendredi: 'Samedi 16 janvier',
    samedi: 'Dimanche 17 janvier',
  } as const,
  dateShort: '16 & 17 janvier 2027',
  dateBadge: '16 & 17 JAN. 2027',
  dateKicker: '16 & 17 JANVIER 2027 À NANTERRE',
  dateOverlay: 'RENDEZ-VOUS LES 16 ET 17 JANVIER À NANTERRE',
  prices: { day: 20, passPresale: 30, passRegular: 35 } as const,
  /** Fin de prévente pass (Europe/Paris) — à partir de cet instant : tarif régulier. */
  passPresaleEndsAt: '2027-01-01T00:00:00+01:00',
  standPrices: { day: 80, pass: 130 } as const,
  kidsFreeUnder: 12,
  currency: 'EUR',
  maxTickets: 20,
  openingTime: '10h30',
  cities: ['Moroni', 'Iconi', 'Itsandra', 'Domoni', 'Mutsamudu', 'Ntsoudjini'] as const,
  venue: {
    name: 'Salle des congrès',
    address: '118 Rue du 8 Mai 1945, 92000 Nanterre',
    phones: ['07 77 16 86 84', '07 71 73 31 13'],
  },
})

export const KIDS_FREE_LABEL = `Gratuit pour les enfants de moins de ${EVENT_CONFIG.kidsFreeUnder} ans`

export const PASS_PRESALE_NOTE =
  `* Tarif prévente jusqu’au 1er janvier 2027 ; ${EVENT_CONFIG.prices.passRegular} € ensuite.`

export type TicketFormula = 'day' | 'pass'
export type TicketDay = keyof typeof EVENT_CONFIG.dates

const passPresaleEndsAt = new Date(EVENT_CONFIG.passPresaleEndsAt)

export function isPassPresale(now: Date = new Date()): boolean {
  return now < passPresaleEndsAt
}

export function getPassPrice(now: Date = new Date()): number {
  return isPassPresale(now) ? EVENT_CONFIG.prices.passPresale : EVENT_CONFIG.prices.passRegular
}

export function getTicketUnitPrice(formula: TicketFormula, now: Date = new Date()): number {
  return formula === 'pass' ? getPassPrice(now) : EVENT_CONFIG.prices.day
}

export function formatEuro(value: number): string {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: EVENT_CONFIG.currency,
    maximumFractionDigits: 0,
  }).format(value)
}
