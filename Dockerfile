# Build stage
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Install migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o news-feed-bot ./cmd/main.go

# Final stage
FROM alpine:3.21

# Install dependencies
RUN apk add --no-cache \
    postgresql15-client \
    tzdata \
    libc6-compat

WORKDIR /app

# Copy necessary files
COPY --from=builder /app/news-feed-bot .
COPY --from=builder /app/wait-for-postgres.sh .
COPY --from=builder /app/.env .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations/schema /app/migrations/schema
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Environment variables
ENV CONFIG_PATH=/app/config/config.yaml

# Set permissions
RUN chmod +x news-feed-bot wait-for-postgres.sh

# Entrypoint
CMD ["./wait-for-postgres.sh", "db:5432", "./news-feed-bot"]