# ─── Build stage ──────────────────────────────────────────────────────
FROM golang:1.24.2-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /music-awards ./main.go

# ─── Runtime stage ────────────────────────────────────────────────────
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app

# copy the compiled binary
COPY --from=builder /music-awards .

# copy your migrations directory too!
COPY --from=builder /app/migrations ./migrations

EXPOSE 8000
STOPSIGNAL SIGINT

CMD ["/app/music-awards"]
