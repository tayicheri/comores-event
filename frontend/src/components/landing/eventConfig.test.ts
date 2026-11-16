import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import {
  EVENT_CONFIG,
  getPassPrice,
  getTicketUnitPrice,
  isPassPresale,
} from './eventConfig.ts'

describe('isPassPresale / getPassPrice', () => {
  it('prévente avant le 1er janvier 2027 00:00 Europe/Paris', () => {
    const before = new Date('2026-12-31T23:59:59+01:00')
    assert.equal(isPassPresale(before), true)
    assert.equal(getPassPrice(before), EVENT_CONFIG.prices.passPresale)
    assert.equal(getTicketUnitPrice('pass', before), 30)
  })

  it('tarif régulier à partir du 1er janvier 2027 00:00', () => {
    const at = new Date(EVENT_CONFIG.passPresaleEndsAt)
    assert.equal(isPassPresale(at), false)
    assert.equal(getPassPrice(at), EVENT_CONFIG.prices.passRegular)
    assert.equal(getTicketUnitPrice('pass', at), 35)
  })

  it('journée inchangée', () => {
    assert.equal(getTicketUnitPrice('day', new Date('2026-06-01')), EVENT_CONFIG.prices.day)
    assert.equal(getTicketUnitPrice('day', new Date('2027-02-01')), EVENT_CONFIG.prices.day)
  })
})
