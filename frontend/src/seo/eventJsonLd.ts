import { EVENT_CONFIG, getPassPrice } from '../components/landing/eventConfig.ts'

/** Aligné sur la meta description de `index.html`. */
export const EVENT_SEO_DESCRIPTION =
  'Les Comores, un patrimoine mondial à célébrer. Salon de la culture et de la gastronomie comorienne — 16 et 17 janvier 2027 à Nanterre.'

export function normalizeSeoSiteUrl(siteUrl: string): string {
  return siteUrl.trim().replace(/\/$/, '')
}

/** JSON-LD schema.org/Event pour la landing (Google + crawlers IA). */
export function buildEventJsonLd(siteUrl: string, now: Date = new Date()): Record<string, unknown> {
  const base = normalizeSeoSiteUrl(siteUrl)
  if (!base) {
    throw new Error('VITE_SITE_URL est requis pour générer le JSON-LD Event')
  }

  const ticketUrl = `${base}/#billetterie`

  return {
    '@context': 'https://schema.org',
    '@type': 'Event',
    name: EVENT_CONFIG.fullName,
    description: EVENT_SEO_DESCRIPTION,
    startDate: EVENT_CONFIG.startDate,
    endDate: EVENT_CONFIG.endDate,
    eventStatus: 'https://schema.org/EventScheduled',
    eventAttendanceMode: 'https://schema.org/OfflineEventAttendanceMode',
    url: `${base}/`,
    image: [`${base}/og-image.jpg`],
    inLanguage: 'fr-FR',
    location: {
      '@type': 'Place',
      name: EVENT_CONFIG.venue.name,
      address: {
        '@type': 'PostalAddress',
        streetAddress: EVENT_CONFIG.venue.streetAddress,
        addressLocality: EVENT_CONFIG.venue.addressLocality,
        postalCode: EVENT_CONFIG.venue.postalCode,
        addressCountry: EVENT_CONFIG.venue.addressCountry,
      },
    },
    offers: [
      {
        '@type': 'Offer',
        name: 'Billet 1 journée',
        price: EVENT_CONFIG.prices.day,
        priceCurrency: EVENT_CONFIG.currency,
        url: ticketUrl,
        availability: 'https://schema.org/InStock',
        category: 'Primary',
      },
      {
        '@type': 'Offer',
        name: 'Pass 2 jours',
        price: getPassPrice(now),
        priceCurrency: EVENT_CONFIG.currency,
        url: ticketUrl,
        availability: 'https://schema.org/InStock',
      },
    ],
    organizer: {
      '@type': 'Organization',
      name: EVENT_CONFIG.fullName,
      url: `${base}/`,
    },
  }
}
