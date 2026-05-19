FROM golang:1.23-alpine AS base
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

# Стадия для миграций
FROM base AS migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.19.1
COPY migrations/ ./migrations/
WORKDIR /app/migrations

# Dev-стадия
FROM base AS dev
COPY . .
WORKDIR /app
CMD ["go", "run", "./cmd/api"]

# Production-стадия
FROM base AS prod
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api
EXPOSE 8080
CMD ["./main"]
