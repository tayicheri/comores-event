/** Contenu robots.txt / sitemap.xml généré au build (SSG). */

export function normalizeSiteUrl(siteUrl: string): string {
  return siteUrl.trim().replace(/\/$/, '')
}

export function generateRobotsTxt(siteUrl: string): string {
  const base = normalizeSiteUrl(siteUrl)
  if (!base) {
    throw new Error('VITE_SITE_URL est requis pour générer robots.txt')
  }
  return [
    'User-agent: *',
    'Allow: /',
    'Disallow: /lekeza',
    '',
    `Sitemap: ${base}/sitemap.xml`,
    '',
  ].join('\n')
}

export function generateSitemapXml(siteUrl: string): string {
  const base = normalizeSiteUrl(siteUrl)
  if (!base) {
    throw new Error('VITE_SITE_URL est requis pour générer sitemap.xml')
  }
  return [
    '<?xml version="1.0" encoding="UTF-8"?>',
    '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
    '  <url>',
    `    <loc>${base}/</loc>`,
    '    <changefreq>weekly</changefreq>',
    '    <priority>1.0</priority>',
    '  </url>',
    '</urlset>',
    '',
  ].join('\n')
}
