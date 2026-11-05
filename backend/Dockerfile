# --- ÉTAPE 1 : BASE ---
FROM golang:1.25-alpine AS base
WORKDIR /app
ENV GOTOOLCHAIN=auto
RUN apk add --no-cache git postgresql-client tzdata

COPY go.mod go.sum ./
RUN go mod download

# --- ÉTAPE 2 : DÉVELOPPEMENT (hot reload via air) ---
FROM base AS dev
RUN go install github.com/air-verse/air@v1.52.3
COPY . .
RUN chmod +x scripts/docker-entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/bin/sh", "/app/scripts/docker-entrypoint.sh"]
CMD ["air", "-c", ".air.toml"]

# --- ÉTAPE 3 : BUILDER ---
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/main ./cmd/api

# --- ÉTAPE 4 : PRODUCTION ---
FROM alpine:3.20 AS prod
RUN apk add --no-cache ca-certificates postgresql-client tzdata \
    && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /out/main /app/main
COPY --from=builder /app/scripts/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
USER app
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/main"]
