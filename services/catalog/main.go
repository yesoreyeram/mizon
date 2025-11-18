package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"mizon/loggerx"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var searchAPI string

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Category    string             `json:"category" bson:"category"`
	Stock       int                `json:"stock" bson:"stock"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

func initMongo() error {
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getEnv("MONGO_DB", "mizon_catalog")
	searchAPI = getEnv("SEARCH_API", "http://localhost:8003")

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
				collection = client.Database(mongoDB).Collection("products")
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

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		http.Error(w, "Error decoding products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product Product
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	product.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}

	product.ID = result.InsertedID.(primitive.ObjectID)

	// Best-effort indexing into search
	go func(p Product) {
		// Build payload expected by search service
		payload := map[string]interface{}{
			"id":          p.ID.Hex(),
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
			"category":    p.Category,
			"stock":       p.Stock,
			"image_url":   p.ImageURL,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			loggerx.Errorf("catalog: error marshaling product for indexing: %v", err)
			return
		}
		url := fmt.Sprintf("%s/api/search/index", searchAPI)
		resp, err := http.Post(url, "application/json", bytes.NewReader(data))
		if err != nil {
			loggerx.Errorf("catalog: failed to index product %s: %v", p.ID.Hex(), err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			loggerx.Warnf("catalog: indexing product %s returned status %d", p.ID.Hex(), resp.StatusCode)
		}
	}(product)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// Reindex all products from catalog into search service
func reindexHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	indexed := 0
	failed := 0
	for cursor.Next(ctx) {
		var p Product
		if err := cursor.Decode(&p); err != nil {
			failed++
			continue
		}
		payload := map[string]interface{}{
			"id":          p.ID.Hex(),
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
			"category":    p.Category,
			"stock":       p.Stock,
			"image_url":   p.ImageURL,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			failed++
			continue
		}
		url := fmt.Sprintf("%s/api/search/index", searchAPI)
		resp, err := http.Post(url, "application/json", bytes.NewReader(data))
		if err != nil {
			failed++
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 300 {
			failed++
			continue
		}
		indexed++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"indexed": indexed, "failed": failed})
}

func getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	categories, err := collection.Distinct(ctx, "category", bson.M{})
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func getProductsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"category": category})
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		http.Error(w, "Error decoding products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
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
	if err := initMongo(); err != nil {
		loggerx.Fatalf("%v", err)
	}

	router := mux.NewRouter()
	cfg := loggerx.Config{LogRequestBody: loggerx.EnvBool("LOG_REQUEST_BODY", false), MaxBody: loggerx.EnvInt("LOG_MAX_BODY", 2048)}
	router.Use(loggerx.Middleware(cfg))
	router.HandleFunc("/api/catalog/products", enableCORS(getProductsHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/catalog/products/{id}", enableCORS(getProductHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/catalog/products", enableCORS(createProductHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/catalog/categories", enableCORS(getCategoriesHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/catalog/categories/{category}/products", enableCORS(getProductsByCategoryHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/catalog/reindex", enableCORS(reindexHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/health", enableCORS(healthHandler)).Methods("GET", "OPTIONS")

	port := getEnv("PORT", "8002")
	loggerx.Infof("Catalog service starting on port %s", port)
	loggerx.Fatalf("%v", http.ListenAndServe(":"+port, router))
}
