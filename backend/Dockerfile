# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application, migrate, and seed binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binaries from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/seed .
COPY --from=builder /app/migrations ./migrations

# Expose the port the app runs on
EXPOSE 8080

# Run migrations and seed before starting the server
CMD sh -c "./migrate && ./seed && ./main"
