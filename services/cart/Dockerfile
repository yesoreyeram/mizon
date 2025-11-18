FROM golang AS builder

# Copy shared logger module and service module files
WORKDIR /app
COPY pkg/loggerx /app/pkg/loggerx

# Set workdir to service and download deps
WORKDIR /app/services/cart
COPY services/cart/go.mod ./go.mod
RUN go mod download

# Copy service sources
COPY services/cart/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o cart-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/cart/cart-service .

EXPOSE 8004

CMD ["./cart-service"]
