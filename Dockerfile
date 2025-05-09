# Build stage
FROM golang:1.24.2-alpine AS builder
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary (using cmd/app/main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -o /music-awards ./cmd/app/main.go

# Runtime stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /

# Copy the compiled binary
COPY --from=builder /music-awards .

# Expose application port
EXPOSE 8000

# Start the server
CMD ["/music-awards"]
