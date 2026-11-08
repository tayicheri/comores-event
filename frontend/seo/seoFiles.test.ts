import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import { generateRobotsTxt, generateSitemapXml, normalizeSiteUrl } from './seoFiles.ts'

describe('normalizeSiteUrl', () => {
  it('retire le slash final', () => {
    assert.equal(normalizeSiteUrl('https://example.com/'), 'https://example.com')
  })
})

describe('generateRobotsTxt', () => {
  it('autorise / et bloque /lekeza', () => {
    const txt = generateRobotsTxt('https://example.com/')
    assert.match(txt, /^User-agent: \*$/m)
    assert.match(txt, /^Allow: \/$/m)
    assert.match(txt, /^Disallow: \/lekeza$/m)
    assert.match(txt, /^Sitemap: https:\/\/example\.com\/sitemap\.xml$/m)
    assert.doesNotMatch(txt, /Allow: \/lekeza/)
  })

  it('exige une URL', () => {
    assert.throws(() => generateRobotsTxt(''), /VITE_SITE_URL/)
  })
})

describe('generateSitemapXml', () => {
  it('n’inclut que la landing', () => {
    const xml = generateSitemapXml('https://example.com')
    assert.match(xml, /<loc>https:\/\/example\.com\/<\/loc>/)
    assert.doesNotMatch(xml, /lekeza/)
  })

  it('exige une URL', () => {
    assert.throws(() => generateSitemapXml('  '), /VITE_SITE_URL/)
  })
})
