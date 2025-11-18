FROM golang AS builder

# Copy shared logger module and service module files
WORKDIR /app
COPY pkg/loggerx /app/pkg/loggerx

# Set workdir to service and download deps
WORKDIR /app/services/search
COPY services/search/go.mod ./go.mod
RUN go mod download

# Copy service sources
COPY services/search/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o search-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/search/search-service .

EXPOSE 8003

CMD ["./search-service"]
