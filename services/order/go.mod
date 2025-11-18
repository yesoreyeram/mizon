module order-service

go 1.23

require (
	github.com/confluentinc/confluent-kafka-go v1.9.2
	github.com/gorilla/mux v1.8.1
	go.mongodb.org/mongo-driver v1.13.1
	mizon/loggerx v0.0.0
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace mizon/loggerx => ../../pkg/loggerx
