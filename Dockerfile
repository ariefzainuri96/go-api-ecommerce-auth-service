# ============================
# 1 — Build Stage
# ============================
FROM golang:1.24 AS builder

WORKDIR /app

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/api

# ============================
# 2 — Runtime Stage
# ============================
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/auth-service .

# Copy migrations folder (this was missing)
COPY --from=builder /app/migrations ./migrations

# Expose the port your service uses
EXPOSE 50051

# Run the application
CMD ["./auth-service"]