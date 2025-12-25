# Gunakan image Golang berbasis Debian (Bookworm) yang lebih stabil
# Ini akan mengatasi error 'pread64' / SQLite yang kamu alami
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Copy file dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy sisa source code
COPY . .

# Build aplikasi
RUN go build -o main .

# --- Stage Akhir (Untuk Menjalankan Aplikasi) ---
FROM debian:bookworm-slim

WORKDIR /app

# Install sertifikat SSL (Wajib supaya bisa connect ke Neon/HTTPS)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy hasil build dari stage sebelumnya
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]