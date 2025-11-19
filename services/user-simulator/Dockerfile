# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy shared packages
COPY pkg/ ./pkg/

# Copy service code
COPY services/user-simulator/ ./services/user-simulator/

# Build the service
WORKDIR /app/services/user-simulator
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-simulator

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /user-simulator .

CMD ["./user-simulator"]
