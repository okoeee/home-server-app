FROM node:24-bookworm-slim AS builder

WORKDIR /app

# 依存定義を先にコピーし、ソースだけの変更ではnpm ciのキャッシュを再利用する。
COPY package.json package-lock.json ./

RUN npm ci

COPY . .

RUN npm run build

FROM caddy:2-alpine

COPY Caddyfile /etc/caddy/Caddyfile

COPY --from=builder /app/dist /usr/share/caddy

EXPOSE 80
