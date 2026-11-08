# syntax=docker/dockerfile:1

# ---------- dev : Vite HMR sur :3000 ----------
FROM node:22-alpine AS dev
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
EXPOSE 3000
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3000"]

# ---------- builder : compile les assets statiques ----------
FROM node:22-alpine AS builder
WORKDIR /app
ARG VITE_API_URL=http://localhost:8080
ARG VITE_APP_NAME="Comores Event"
ARG VITE_SITE_URL=
ENV VITE_API_URL=$VITE_API_URL \
    VITE_APP_NAME=$VITE_APP_NAME \
    VITE_SITE_URL=$VITE_SITE_URL
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
RUN npm run build

# ---------- prod : nginx alpine ----------
FROM nginx:1.27-alpine AS prod
ARG VITE_API_URL=http://localhost:8080
# Origine de l'API autorisée dans la CSP (connect-src). Surchargeable au runtime.
ENV API_ORIGIN=$VITE_API_URL
# L'entrypoint officiel nginx applique envsubst sur /etc/nginx/templates/*.template
COPY nginx.conf.template /etc/nginx/templates/default.conf.template
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD wget -qO- http://127.0.0.1/healthz >/dev/null 2>&1 || exit 1
