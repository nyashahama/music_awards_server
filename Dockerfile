# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /music-awards ./cmd/server

# Run stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /music-awards .
EXPOSE 8080
CMD ["/music-awards"]
