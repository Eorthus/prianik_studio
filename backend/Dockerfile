# Этап сборки
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

# Копируем исходный код
COPY . .

# Компилируем приложение с отключенным CGO и оптимизацией размера
RUN CGO_ENABLED=0 GOOS=linux GOGC=off go build -ldflags="-s -w" -o server ./cmd/server

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из образа builder
COPY --from=builder /app/server .
COPY --from=builder /app/.env .

# Открываем порт
EXPOSE 8080

# Запускаем сервер
CMD ["./server"]