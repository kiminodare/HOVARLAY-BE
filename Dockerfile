# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install git untuk go mod
RUN apk add --no-cache git && \
    go mod download

COPY go.mod go.sum ./
COPY . .

# Build server and migrate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate

# Stage 2: Final
FROM alpine:3.22.1
RUN apk --no-cache add bash ca-certificates tini && \
    addgroup -S app && \
    adduser -S app -G app

USER app

WORKDIR /home/app
COPY --from=builder /app/server .
COPY --from=builder /app/migrate .

EXPOSE 9888

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:9888/health || exit 1

ENTRYPOINT ["/sbin/tini", "--"]

# Jalankan migrate dulu, baru server
CMD ["sh", "-c", "./migrate && ./server"]