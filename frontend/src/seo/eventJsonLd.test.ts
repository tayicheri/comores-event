import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import { EVENT_CONFIG } from '../components/landing/eventConfig.ts'
import { buildEventJsonLd, EVENT_SEO_DESCRIPTION, normalizeSeoSiteUrl } from './eventJsonLd.ts'

describe('normalizeSeoSiteUrl', () => {
  it('retire le slash final', () => {
    assert.equal(normalizeSeoSiteUrl('https://example.com/'), 'https://example.com')
  })
})

describe('buildEventJsonLd', () => {
  const nowPresale = new Date('2026-06-01T12:00:00+02:00')
  const nowRegular = new Date('2027-01-02T12:00:00+01:00')

  it('exige une URL', () => {
    assert.throws(() => buildEventJsonLd(''), /VITE_SITE_URL/)
  })

  it('produit un Event schema.org complet', () => {
    const ld = buildEventJsonLd('https://example.com/', nowPresale)
    assert.equal(ld['@context'], 'https://schema.org')
    assert.equal(ld['@type'], 'Event')
    assert.equal(ld.name, EVENT_CONFIG.fullName)
    assert.equal(ld.description, EVENT_SEO_DESCRIPTION)
    assert.equal(ld.startDate, EVENT_CONFIG.startDate)
    assert.equal(ld.endDate, EVENT_CONFIG.endDate)
    assert.equal(ld.url, 'https://example.com/')
    assert.deepEqual(ld.image, ['https://example.com/og-image.jpg'])
    assert.equal(ld.eventStatus, 'https://schema.org/EventScheduled')
    assert.equal(ld.eventAttendanceMode, 'https://schema.org/OfflineEventAttendanceMode')

    const location = ld.location as Record<string, unknown>
    assert.equal(location['@type'], 'Place')
    assert.equal(location.name, EVENT_CONFIG.venue.name)
    const address = location.address as Record<string, unknown>
    assert.equal(address.streetAddress, EVENT_CONFIG.venue.streetAddress)
    assert.equal(address.addressLocality, 'Nanterre')
    assert.equal(address.postalCode, '92000')
    assert.equal(address.addressCountry, 'FR')

    const offers = ld.offers as Record<string, unknown>[]
    assert.equal(offers.length, 2)
    assert.equal(offers[0].price, EVENT_CONFIG.prices.day)
    assert.equal(offers[0].priceCurrency, 'EUR')
    assert.equal(offers[0].url, 'https://example.com/#billetterie')
    assert.equal(offers[1].price, EVENT_CONFIG.prices.passPresale)

    const organizer = ld.organizer as Record<string, unknown>
    assert.equal(organizer['@type'], 'Organization')
    assert.equal(organizer.url, 'https://example.com/')
  })

  it('utilise le tarif pass régulier après la prévente', () => {
    const ld = buildEventJsonLd('https://example.com', nowRegular)
    const offers = ld.offers as Record<string, unknown>[]
    assert.equal(offers[1].price, EVENT_CONFIG.prices.passRegular)
  })
})
