# Multi-stage Dockerfile for Go Authentication Module

# Stage 1: Build
FROM golang:1.20-slim AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main ./cmd/main.go

# Stage 2: Production
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 80

CMD ["./main"]
