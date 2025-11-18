package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"mizon/loggerx"
	"mizon/telemetry"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
)

var esClient *elasticsearch.Client

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
}

type SearchResult struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
}

func initElasticsearch() error {
	esURL := getEnv("ELASTICSEARCH_URL", "http://localhost:9200")

	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
		Transport: telemetry.OtelTransport(),
	}

	var err error
	for i := 0; i < 30; i++ {
		esClient, err = elasticsearch.NewClient(cfg)
		if err == nil {
			_, err = esClient.Ping()
			if err == nil {
				loggerx.Info("Successfully connected to Elasticsearch")
				return createIndex()
			}
		}
		loggerx.Infof("Waiting for Elasticsearch... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("failed to connect to Elasticsearch: %v", err)
}

func createIndex() error {
	indexName := "products"

	mapping := `{
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"name": { "type": "text" },
				"description": { "type": "text" },
				"price": { "type": "float" },
				"category": { "type": "keyword" },
				"stock": { "type": "integer" },
				"image_url": { "type": "keyword" }
			}
		}
	}`

	res, err := esClient.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		res, err = esClient.Indices.Create(
			indexName,
			esClient.Indices.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		loggerx.Info("Created products index")
	}

	return nil
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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	searchQuery := fmt.Sprintf(`{
		"query": {
			"multi_match": {
				"query": "%s",
				"fields": ["name^2", "description", "category"]
			}
		},
		"size": 50
	}`, query)

	res, err := esClient.Search(
		esClient.Search.WithContext(r.Context()),
		esClient.Search.WithIndex("products"),
		esClient.Search.WithBody(strings.NewReader(searchQuery)),
	)
	if err != nil {
		loggerx.Errorf("Error searching: %v", err)
		http.Error(w, "Search error", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		http.Error(w, "Search error", http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		http.Error(w, "Error parsing response", http.StatusInternalServerError)
		return
	}

	hits := result["hits"].(map[string]interface{})
	total := int(hits["total"].(map[string]interface{})["value"].(float64))
	hitList := hits["hits"].([]interface{})

	var products []Product
	for _, hit := range hitList {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		product := Product{
			ID:          source["id"].(string),
			Name:        source["name"].(string),
			Description: source["description"].(string),
			Price:       source["price"].(float64),
			Category:    source["category"].(string),
			Stock:       int(source["stock"].(float64)),
			ImageURL:    source["image_url"].(string),
		}
		products = append(products, product)
	}

	searchResult := SearchResult{
		Products: products,
		Total:    total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResult)
}

func indexProductHandler(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(product)
	if err != nil {
		http.Error(w, "Error encoding product", http.StatusInternalServerError)
		return
	}

	res, err := esClient.Index(
		"products",
		bytes.NewReader(data),
		esClient.Index.WithDocumentID(product.ID),
		esClient.Index.WithRefresh("true"),
	)
	if err != nil {
		loggerx.Errorf("Error indexing product: %v", err)
		http.Error(w, "Error indexing product", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		http.Error(w, "Error indexing product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "indexed", "id": product.ID})
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
	if _, err := telemetry.Setup("search-service"); err != nil {
		loggerx.Warnf("tracing setup failed: %v", err)
	}
	if err := initElasticsearch(); err != nil {
		loggerx.Fatalf("%v", err)
	}

	router := mux.NewRouter()
	router.Use(telemetry.MuxMiddleware("search-service"))
	cfg := loggerx.Config{LogRequestBody: loggerx.EnvBool("LOG_REQUEST_BODY", false), MaxBody: loggerx.EnvInt("LOG_MAX_BODY", 2048)}
	router.Use(loggerx.Middleware(cfg))
	router.HandleFunc("/api/search", enableCORS(searchHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/search/index", enableCORS(indexProductHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/health", enableCORS(healthHandler)).Methods("GET", "OPTIONS")

	port := getEnv("PORT", "8003")
	loggerx.Infof("Search service starting on port %s", port)
	loggerx.Fatalf("%v", http.ListenAndServe(":"+port, router))
}
