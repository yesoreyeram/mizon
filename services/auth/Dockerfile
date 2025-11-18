FROM golang AS builder

# Copy shared logger module and service module files
WORKDIR /app
COPY pkg/loggerx /app/pkg/loggerx

# Set workdir to service and download deps
WORKDIR /app/services/auth
COPY services/auth/go.mod ./go.mod
RUN go mod download

# Copy service sources
COPY services/auth/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/auth/auth-service .

EXPOSE 8001

CMD ["./auth-service"]
