# Step 1: Build stage
FROM golang:1.23 AS builder

# Set working dir
WORKDIR /app

# Copy go.mod and go.sum first (for dependency cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy full source
COPY . .

# Build static binary (no CGO, smaller runtime deps)
RUN CGO_ENABLED=0 GOOS=linux go build -o go-service ./cmd/main.go

# Step 2: Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

# Install CA certs (needed for HTTPS/TLS connections)
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /app/go-service .

# Expose port
EXPOSE 9090

# Run app
CMD ["./go-service"]
