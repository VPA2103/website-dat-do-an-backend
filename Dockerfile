za# ==========================
# Stage 1: Build Go App
# ==========================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod trước để cache
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source
COPY . .

# Build file thực thi
RUN go build -o app .

# ==========================
# Stage 2: Run Application
# ==========================
FROM alpine:latest

WORKDIR /app

# Cài chứng chỉ & timezone nếu app gọi HTTPS (Cloudinary cần)
RUN apk add --no-cache ca-certificates tzdata

# Copy binary từ builder
COPY --from=builder /app/app .


EXPOSE 30000

CMD ["./app"]
