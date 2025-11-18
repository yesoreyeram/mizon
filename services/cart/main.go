package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"mizon/loggerx"
	"mizon/telemetry"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

type CartItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Name      string  `json:"name" bson:"name"`
	Price     float64 `json:"price" bson:"price"`
	Quantity  int     `json:"quantity" bson:"quantity"`
}

type Cart struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Items     []CartItem         `json:"items" bson:"items"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type AddItemRequest struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Name      string  `json:"name" bson:"name"`
	Price     float64 `json:"price" bson:"price"`
	Quantity  int     `json:"quantity" bson:"quantity"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity"`
}

func initMongo() error {
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getEnv("MONGO_DB", "mizon_cart")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client *mongo.Client
	var err error

	for i := 0; i < 30; i++ {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI).SetMonitor(telemetry.MongoMonitor()))
		if err == nil {
			err = client.Ping(ctx, nil)
			if err == nil {
				loggerx.Info("Successfully connected to MongoDB")
				collection = client.Database(mongoDB).Collection("carts")
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

func getCartHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var cart Cart
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return empty cart
			cart = Cart{
				UserID:    userID,
				Items:     []CartItem{},
				UpdatedAt: time.Now(),
			}
		} else {
			http.Error(w, "Error fetching cart", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func addItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	newItem := CartItem(req)

	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set":         bson.M{"updated_at": time.Now()},
		"$setOnInsert": bson.M{"user_id": userID},
	}

	// Check if item already exists
	var cart Cart
	err := collection.FindOne(ctx, filter).Decode(&cart)

	if err == mongo.ErrNoDocuments {
		// Create new cart
		update["$set"].(bson.M)["items"] = []CartItem{newItem}
	} else if err == nil {
		// Update existing cart
		found := false
		for i, item := range cart.Items {
			if item.ProductID == req.ProductID {
				cart.Items[i].Quantity += req.Quantity
				found = true
				break
			}
		}
		if !found {
			cart.Items = append(cart.Items, newItem)
		}
		update["$set"].(bson.M)["items"] = cart.Items
	} else {
		http.Error(w, "Error accessing cart", http.StatusInternalServerError)
		return
	}

	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		http.Error(w, "Error updating cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "added"})
}

func updateItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	itemID := vars["itemId"]

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var cart Cart
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	for i, item := range cart.Items {
		if item.ProductID == itemID {
			if req.Quantity <= 0 {
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				cart.Items[i].Quantity = req.Quantity
			}
			break
		}
	}

	update := bson.M{
		"$set": bson.M{
			"items":      cart.Items,
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	if err != nil {
		http.Error(w, "Error updating cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func removeItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	itemID := vars["itemId"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var cart Cart
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	for i, item := range cart.Items {
		if item.ProductID == itemID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			break
		}
	}

	update := bson.M{
		"$set": bson.M{
			"items":      cart.Items,
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	if err != nil {
		http.Error(w, "Error updating cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "removed"})
}

func clearCartHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"items":      []CartItem{},
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	if err != nil {
		http.Error(w, "Error clearing cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "cleared"})
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
	if _, err := telemetry.Setup("cart-service"); err != nil {
		loggerx.Warnf("tracing setup failed: %v", err)
	}
	if err := initMongo(); err != nil {
		loggerx.Fatalf("%v", err)
	}

	router := mux.NewRouter()
	router.Use(telemetry.MuxMiddleware("cart-service"))
	cfg := loggerx.Config{LogRequestBody: loggerx.EnvBool("LOG_REQUEST_BODY", false), MaxBody: loggerx.EnvInt("LOG_MAX_BODY", 2048)}
	router.Use(loggerx.Middleware(cfg))
	router.HandleFunc("/api/cart/{userId}", enableCORS(getCartHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/cart/{userId}/items", enableCORS(addItemHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/cart/{userId}/items/{itemId}", enableCORS(updateItemHandler)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/cart/{userId}/items/{itemId}", enableCORS(removeItemHandler)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/cart/{userId}", enableCORS(clearCartHandler)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/health", enableCORS(healthHandler)).Methods("GET", "OPTIONS")

	port := getEnv("PORT", "8004")
	loggerx.Infof("Cart service starting on port %s", port)
	loggerx.Fatalf("%v", http.ListenAndServe(":"+port, router))
}
