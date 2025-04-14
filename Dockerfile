FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o news-feed-bot ./cmd/main.go

FROM alpine:3.21

RUN apk add --no-cache \
    postgresql15-client \
    tzdata \
    libc6-compat

WORKDIR /app

COPY --from=builder /app/news-feed-bot .
COPY --from=builder /app/wait-for-postgres.sh .
COPY --from=builder /app/.env .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations/schema /app/migrations/schema
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

ENV CONFIG_PATH=/app/config/config.yaml

RUN chmod +x news-feed-bot wait-for-postgres.sh

CMD ["./wait-for-postgres.sh", "db", "5432", "./news-feed-bot"]