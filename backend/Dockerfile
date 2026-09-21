# Стадия сборки
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Финальный образ — минимальный, без Go-тулчейна
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

COPY --from=builder /server /server

EXPOSE 3000

CMD ["/server"]
