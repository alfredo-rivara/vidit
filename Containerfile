# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary suitable for alpine/scratch
RUN CGO_ENABLED=0 GOOS=linux go build -o vidit-server ./cmd/server

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/vidit-server .
# Copy css and templates as they are read at runtime
COPY --from=builder /app/public ./public
COPY --from=builder /app/views ./views

# Set timezone (Optional, good for news)
ENV TZ=America/Santiago

# Expose port
EXPOSE 3000

# Run
CMD ["./vidit-server"]
