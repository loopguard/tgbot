FROM golang:1.21-alpine as builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tgbot cmd/main.go

FROM scratch

COPY --from=builder /tgbot /

ENTRYPOINT ["/tgbot"]