FROM golang AS builder

# Copy shared logger module and service module files
WORKDIR /app
COPY pkg/loggerx /app/pkg/loggerx

# Set workdir to service and download deps
WORKDIR /app/services/catalog
COPY services/catalog/go.mod ./go.mod
RUN go mod download

# Copy service sources
COPY services/catalog/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o catalog-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/catalog/catalog-service .

EXPOSE 8002

CMD ["./catalog-service"]
