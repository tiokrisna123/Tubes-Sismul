FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example ./.env

# Create data directory for SQLite
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./main"]
