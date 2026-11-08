import { writeFileSync } from 'node:fs'
import { join } from 'node:path'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { generateRobotsTxt, generateSitemapXml } from './seo/seoFiles'

function writeSeoFiles(distDir: string, siteUrl: string) {
  writeFileSync(join(distDir, 'robots.txt'), generateRobotsTxt(siteUrl), 'utf8')
  writeFileSync(join(distDir, 'sitemap.xml'), generateSitemapXml(siteUrl), 'utf8')
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const devProxyTarget = env.VITE_DEV_PROXY_API || process.env.VITE_DEV_PROXY_API
  const siteUrl = (env.VITE_SITE_URL || process.env.VITE_SITE_URL || '').replace(/\/$/, '')
  const ogUrl = siteUrl || '/'
  const ogImage = siteUrl ? `${siteUrl}/og-image.jpg` : '/og-image.jpg'
  const distDir = fileURLToPath(new URL('./dist', import.meta.url))

  return {
    plugins: [
      vue(),
      tailwindcss(),
      {
        name: 'inject-open-graph-urls',
        transformIndexHtml(html) {
          return html.replaceAll('__OG_URL__', ogUrl).replaceAll('__OG_IMAGE__', ogImage)
        },
      },
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      host: '0.0.0.0',
      port: 3000,
      strictPort: true,
      watch: { usePolling: true },
      proxy: devProxyTarget
        ? {
            '/api': { target: devProxyTarget, changeOrigin: true },
            '/health': { target: devProxyTarget, changeOrigin: true },
          }
        : undefined,
    },
    preview: {
      host: '0.0.0.0',
      port: 3000,
    },
    ssgOptions: {
      includedRoutes() {
        return ['/']
      },
      onFinished() {
        if (!siteUrl) {
          throw new Error('VITE_SITE_URL est requis pour le build SSG (robots.txt / sitemap.xml)')
        }
        writeSeoFiles(distDir, siteUrl)
      },
    },
  }
})
