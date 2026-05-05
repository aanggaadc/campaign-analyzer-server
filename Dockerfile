# ========= STAGE 1: BUILD =========
FROM golang:1.26.2-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependency untuk build (optional tapi aman)
RUN apk add --no-cache git

# Copy go.mod dan go.sum dulu (cache dependency)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

# ========= STAGE 2: RUN =========
FROM alpine:latest

# Add CA certificates (untuk HTTPS call dari backend)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary dari builder
COPY --from=builder /app/main .

# Expose port (sesuaikan dengan app kamu)
EXPOSE 8080

# Jalankan app
CMD ["./main"]