# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install git (jika diperlukan) dan dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Stage 2: Final
FROM alpine:latest
RUN apk --no-cache add ca-certificates tini

# Tambah user non-root
RUN addgroup -S app && adduser -S app -G app
USER app

WORKDIR /home/app
COPY --from=builder /app/main .

EXPOSE 9888

# Healthcheck (cek endpoint /health misal)
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:9888/health || exit 1

# Gunakan tini sebagai init process
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["./main"]
