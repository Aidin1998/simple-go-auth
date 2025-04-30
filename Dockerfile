# ─── Build Stage ────────────────────────────────────────────────────────────────
FROM golang:1.24 AS builder

# Create working directory
WORKDIR /workspace

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code & build the CLI
COPY . .

# Compile the auth server from cmd/main.go
# Produces a statically linked binary named simple-go-auth
RUN CGO_ENABLED=0 GOOS=linux go build -o /workspace/simple-go-auth ./cmd

# ─── Final Stage ────────────────────────────────────────────────────────────────
FROM alpine:3.18

# Install ca-certificates for HTTPS support
RUN apk add --no-cache ca-certificates

# Copy the compiled binary
COPY --from=builder /workspace/simple-go-auth /usr/local/bin/simple-go-auth

# Expose port 80 to match your ALB listener
EXPOSE 80

# Run the server
ENTRYPOINT ["/usr/local/bin/simple-go-auth"]
