# Frontend SSG & SEO

La landing (`/`) est **prérendue au build** (SSG via `vite-ssg`). Admin et login restent du SPA hydraté côté client. La prod sert toujours des fichiers statiques via nginx (pas de Node SSR runtime).

## Build

```bash
# dans le conteneur web / stage builder
npm run build   # vue-tsc + vite-ssg build
```

Deps build : `vite-ssg` + `@unhead/vue@2` (peer ; ne pas monter en v3 — incompatible avec le head SSR de vite-ssg).

Variables build importantes :

| Variable | Rôle |
|----------|------|
| `VITE_SITE_URL` | URL canonique (prod : `FRONTEND_URL`) — **obligatoire** pour `robots.txt` / `sitemap.xml` |
| `VITE_API_URL` | Origine API (CSP `connect-src`) |
| `VITE_APP_NAME` | Suffixe titres |

Route prérendue : **`/` uniquement** (`ssgOptions.includedRoutes`).

## Fichiers SEO générés

Écrits dans `dist/` en fin de build (`ssgOptions.onFinished`) :

- `robots.txt` — `Allow: /`, `Disallow: /lekeza`, lien Sitemap
- `sitemap.xml` — uniquement `{VITE_SITE_URL}/`

Test unitaire du helper : `npm run test:seo`.

## Exclusion `/lekeza`

Défense en profondeur (pas un contrôle d’accès — l’auth JWT reste obligatoire) :

1. Absent du sitemap
2. `Disallow: /lekeza` dans `robots.txt` (couvre `/lekeza/login`)
3. Meta `robots: noindex, nofollow` via `@unhead/vue` (`route.meta.noindex`)
4. Header HTTP `X-Robots-Tag: noindex, nofollow, noarchive, nosnippet` — nginx `location ^~ /lekeza`

## Caddy

Le domaine frontend **n’envoie plus** de `X-Robots-Tag: noindex` global (sinon la landing resterait non indexable). L’API conserve le noindex.

## Dev

`npm run dev` reste Vite HMR classique (pas de prérendu). Les guards auth ne s’exécutent pas pendant le prérendu SSG (`isClient` only).
