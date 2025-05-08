# Build stage
FROM golang:1.24.2-alpine AS builder
WORKDIR /app

# Download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /music-awards ./cmd/server

# Runtime stage
FROM alpine:3.18

# Install certificates for HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /

# Copy the compiled binary
COPY --from=builder /music-awards .

# Expose application port
EXPOSE 8080

# Start the server
CMD ["/music-awards"]
