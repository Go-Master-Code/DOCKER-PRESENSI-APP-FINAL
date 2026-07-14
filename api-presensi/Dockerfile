# ==========================
# Stage 1 - Build
# ==========================
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ==========================
# Stage 2 - Runtime
# ==========================
FROM alpine:latest

# - tzdata       : agar timezone Asia/Jakarta tersedia
# - mysql-client : untuk proses backup database

RUN apk add --no-cache \
    mysql-client \
    tzdata

# Set timezone container menjadi WIB
ENV TZ=Asia/Jakarta
# ENV diperlukan karena dari go generate UTC sementara kita mau WIB (Asia/Jakarta)

# Folder kerja aplikasi
WORKDIR /app

# Salin binary hasil build
COPY --from=builder /app/main .

# Buat folder yang akan dipakai sebagai Docker Volume
# Folder ini akan dipersist oleh docker-compose
# -------------------------------------------------------
RUN mkdir -p \
    /app/uploads \
    /app/backups

# Port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]