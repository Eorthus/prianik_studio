# Этап сборки
FROM node:20-alpine AS builder

# Устанавливаем pnpm
RUN corepack enable
RUN corepack prepare pnpm@latest --activate

WORKDIR /app

# Копируем только package.json и pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./
RUN pnpm install

# Копируем исходный код
COPY . .

# Полная очистка кеша и временных файлов
RUN rm -rf .nuxt .output node_modules/.cache
RUN pnpm install --force

# Собираем приложение
ENV NODE_ENV=production
ENV NODE_OPTIONS="--max-old-space-size=8192"
ENV NITRO_PRESET=node-server

RUN npx nuxi build --standalone

# Финальный образ
FROM node:20-alpine

# Устанавливаем pnpm
RUN corepack enable
RUN corepack prepare pnpm@latest --activate

WORKDIR /app

# Копируем сборку
COPY --from=builder /app/.output ./.output
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/pnpm-lock.yaml ./pnpm-lock.yaml

# Устанавливаем только production зависимости
RUN pnpm install --prod

# Запускаем приложение
ENV NODE_ENV=production
EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]