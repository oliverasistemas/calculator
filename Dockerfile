# --- Backend build ---
FROM golang:1.26-alpine AS api-build
WORKDIR /src
COPY services/api/go.mod services/api/go.sum ./services/api/
RUN cd services/api && go mod download
COPY services/ services/
RUN cd services/api && CGO_ENABLED=0 go build -o /api ./cmd

# --- Frontend build ---
FROM node:22-alpine AS web-build
WORKDIR /src
COPY app-web/package.json app-web/package-lock.json ./
RUN npm ci
COPY app-web/ .
ARG VITE_API_URL=/api
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

# --- Runtime ---
FROM alpine:3.20
RUN apk add --no-cache caddy

COPY --from=api-build /api /usr/local/bin/api
COPY --from=web-build /src/dist /srv/web

COPY <<'EOF' /etc/caddy/Caddyfile
:8080 {
    handle /api/* {
        uri strip_prefix /api
        reverse_proxy localhost:4000
    }
    handle {
        root * /srv/web
        try_files {path} /index.html
        file_server
    }
}
EOF

EXPOSE 8080

CMD sh -c '/usr/local/bin/api & caddy run --config /etc/caddy/Caddyfile'
