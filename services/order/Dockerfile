FROM golang:1.25-bookworm AS builder

# Install build dependencies for librdkafka (CGO)
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    build-essential \
    pkg-config \
    librdkafka-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Include shared logger module
COPY pkg/loggerx /app/pkg/loggerx

# Build the service
WORKDIR /app/services/order
COPY services/order/go.mod ./go.mod
RUN go mod download

COPY services/order/ .
RUN go build -o order-service .

# Runtime image (Debian slim) with runtime libs
FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    ca-certificates \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/services/order/order-service .

EXPOSE 8005

CMD ["./order-service"]
