# Этап сборки
FROM golang:1.23-bookworm AS builder

# Устанавливаем зависимости для сборки
RUN apt-get update && apt-get install -y git make

# Создаем рабочую директорию /center
WORKDIR /center

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o monitoring-center ./cmd/app/main.go 

# Финальный этап
FROM debian:bookworm-slim

# Устанавливаем необходимые пакеты
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# Создаем рабочую директорию /center
WORKDIR /center

# Копируем бинарник из этапа сборки
COPY --from=builder /center/monitoring-center .

# Копируем конфигурационные файлы
COPY config/config.yml ./config/config.yml
COPY build/init ./init

# Открываем порт сервера
EXPOSE 8080

# Команда запуска приложения
CMD ["./monitoring-center"]