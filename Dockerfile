# =========================
# Stage 1: Build
# =========================
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Copy dependency files terlebih dahulu
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build aplikasi
RUN go build -o server .

# =========================
# Stage 2: Production
# =========================
FROM alpine:3.22

WORKDIR /app

# Sertifikat diperlukan untuk koneksi HTTPS/TLS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary dari builder
COPY --from=builder /app/server .

# Dokumentasi port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./server"]