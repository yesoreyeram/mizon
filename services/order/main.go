package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"mizon/loggerx"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	producer   *kafka.Producer
	collection *mongo.Collection
)

type OrderItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Name      string  `json:"name" bson:"name"`
	Price     float64 `json:"price" bson:"price"`
	Quantity  int     `json:"quantity" bson:"quantity"`
}

type Order struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Items     []OrderItem        `json:"items" bson:"items"`
	Total     float64            `json:"total" bson:"total"`
	Status    string             `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type CreateOrderRequest struct {
	UserID string      `json:"user_id"`
	Items  []OrderItem `json:"items"`
	Total  float64     `json:"total"`
}

func initKafka() error {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")

	var err error
	for i := 0; i < 30; i++ {
		producer, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": broker,
		})
		if err == nil {
			loggerx.Info("Successfully connected to Kafka")
			go handleKafkaEvents()
			return nil
		}
		loggerx.Infof("Waiting for Kafka... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("failed to connect to Kafka: %v", err)
}

func handleKafkaEvents() {
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				loggerx.Errorf("Failed to deliver message: %v", ev.TopicPartition.Error)
			} else {
				loggerx.Infof("Message delivered to %v", ev.TopicPartition)
			}
		}
	}
}

func initMongo() error {
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getEnv("MONGO_DB", "mizon_orders")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client *mongo.Client
	var err error

	for i := 0; i < 30; i++ {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err == nil {
			err = client.Ping(ctx, nil)
			if err == nil {
				loggerx.Info("Successfully connected to MongoDB")
				collection = client.Database(mongoDB).Collection("orders")
				return nil
			}
		}
		loggerx.Infof("Waiting for MongoDB... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("failed to connect to MongoDB: %v", err)
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	order := Order{
		UserID:    req.UserID,
		Items:     req.Items,
		Total:     req.Total,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	order.ID = result.InsertedID.(primitive.ObjectID)

	// Publish to Kafka
	orderData, _ := json.Marshal(order)
	topic := "orders"
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          orderData,
	}, nil)

	if err != nil {
		loggerx.Errorf("Failed to produce message: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func getUserOrdersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		http.Error(w, "Error fetching orders", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var orders []Order
	if err = cursor.All(ctx, &orders); err != nil {
		http.Error(w, "Error decoding orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func getOrderDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order Order
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&order)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	loggerx.Setup()
	if err := initKafka(); err != nil {
		loggerx.Fatalf("%v", err)
	}
	defer producer.Close()

	if err := initMongo(); err != nil {
		loggerx.Fatalf("%v", err)
	}

	router := mux.NewRouter()
	cfg := loggerx.Config{LogRequestBody: loggerx.EnvBool("LOG_REQUEST_BODY", false), MaxBody: loggerx.EnvInt("LOG_MAX_BODY", 2048)}
	router.Use(loggerx.Middleware(cfg))
	router.HandleFunc("/api/orders", enableCORS(createOrderHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/orders/{userId}", enableCORS(getUserOrdersHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/orders/details/{orderId}", enableCORS(getOrderDetailsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/health", enableCORS(healthHandler)).Methods("GET", "OPTIONS")

	port := getEnv("PORT", "8005")
	loggerx.Infof("Order service starting on port %s", port)
	loggerx.Fatalf("%v", http.ListenAndServe(":"+port, router))
}
