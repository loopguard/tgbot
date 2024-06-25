FROM golang:latest

WORKDIR /
COPY . .

# Собираем приложение
RUN go build -o tgbot cmd/main.go

# Определяем порт, на котором будет работать приложение
EXPOSE 8080

# Запускаем приложение при старте контейнера
CMD ["./tgbot"]