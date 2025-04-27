# Multi-stage Dockerfile for Go Authentication Module

# Stage 1: Build
FROM golang:1.22-alpine AS builder


WORKDIR /app

# Ensure go.mod and go.sum are copied correctly
#COPY go.mod go.sum /app/

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# Ensure the Go version is correct and compatible
RUN go version && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go


# Stage 2: Production
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 80

CMD ["./main"]
