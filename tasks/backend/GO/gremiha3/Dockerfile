# Используем официальный образ Golang для сборки
FROM golang:1.21.7-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь проект в контейнер
COPY . .

# Сборка приложения
RUN go build -o myapp ./cmd/myapp/main.go

# Создаем финальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение из предыдущего контейнера
COPY --from=builder /app/myapp .

# Прокидываем аргументы сборки
ARG JWT_KEY

# Копируем файл конфигурации
COPY ./configs/config.yaml ./configs/
COPY ./bin/key.pem .
COPY ./bin/cert.pem .


# Прокидываем порт 8443
EXPOSE 8443

# Устанавливаем переменные окружения
ENV TLS_KEY="/app/key.pem"
ENV TLS_CERT="/app/cert.pem"
ENV JWT_KEY=${JWT_KEY}


# Запуск приложения
CMD ["./myapp"]
