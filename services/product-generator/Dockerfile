FROM golang AS builder

WORKDIR /app

COPY services/product-generator/go.mod services/product-generator/go.sum ./services/product-generator/
COPY pkg/ ./pkg/

WORKDIR /app/services/product-generator
RUN go mod download

WORKDIR /app
COPY services/product-generator/ ./services/product-generator/

WORKDIR /app/services/product-generator
RUN CGO_ENABLED=0 GOOS=linux go build -o product-generator .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/product-generator/product-generator .

CMD ["./product-generator"]
