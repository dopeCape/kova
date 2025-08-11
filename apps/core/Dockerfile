# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY apps/core/go.mod apps/core/go.sum ./
RUN go mod download

# Copy source code
COPY apps/core/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kova-api ./cmd/api

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl postgresql-client

WORKDIR /app

# Create non-root user
RUN addgroup -g 1000 kova && \
    adduser -D -s /bin/sh -u 1000 -G kova kova

# Copy binary from builder
COPY --from=builder /app/kova-api .
COPY --from=builder /app/internal/store/migrations ./internal/store/migrations

# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Change ownership
RUN chown -R kova:kova /app

USER kova

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["./kova-api"]
