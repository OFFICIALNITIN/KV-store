FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o kv-server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kv-server .
COPY --from=builder /app/config.yaml .

EXPOSE 6379

CMD ["./kv-server"]