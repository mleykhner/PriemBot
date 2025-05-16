# Этап сборки
FROM golang:alpine AS builder

# Установка необходимых зависимостей для сборки
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование файлов зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o priembot .

# Этап финального образа
FROM scratch

# Копирование корневых сертификатов
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копирование бинарного файла из этапа сборки
COPY --from=builder /app/priembot /priembot

# Установка рабочей директории
WORKDIR /

# Запуск приложения
ENTRYPOINT ["/priembot"] 