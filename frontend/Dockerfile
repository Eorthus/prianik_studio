FROM node:20-alpine

# Устанавливаем pnpm
RUN corepack enable
RUN corepack prepare pnpm@latest --activate

WORKDIR /app

# Копируем package.json и pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./
RUN pnpm install

# Копируем исходный код
COPY . .

# Настраиваем переменные окружения
ENV NODE_ENV=development

# Запускаем приложение в режиме разработки
CMD ["pnpm", "run", "dev"]