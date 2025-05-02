FROM golang:1.23-alpine

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

# Копируем исходный код
COPY . .

# Настраиваем переменные окружения
ENV GIN_MODE=debug

# Запускаем приложение из правильной директории
CMD ["go", "run", "./cmd/server/main.go"]