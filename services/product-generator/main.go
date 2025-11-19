package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"mizon/loggerx"
)

type Product struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
}

var (
	categories = []string{
		"Electronics",
		"Home & Kitchen",
		"Sports & Outdoors",
		"Books",
		"Clothing",
		"Toys & Games",
		"Health & Beauty",
		"Automotive",
		"Garden & Tools",
		"Pet Supplies",
	}

	productTypes = []string{
		"Smart", "Wireless", "Premium", "Eco-Friendly", "Portable",
		"Professional", "Compact", "Deluxe", "Classic", "Modern",
		"Vintage", "Digital", "Automatic", "Manual", "Electric",
		"Rechargeable", "Heavy-Duty", "Lightweight", "Durable", "Flexible",
	}

	productNames = []string{
		"Speaker", "Headphones", "Camera", "Watch", "Keyboard",
		"Mouse", "Charger", "Cable", "Stand", "Holder",
		"Case", "Cover", "Protector", "Cleaner", "Organizer",
		"Lamp", "Fan", "Heater", "Cooler", "Humidifier",
		"Blender", "Mixer", "Grinder", "Toaster", "Kettle",
		"Backpack", "Wallet", "Bag", "Bottle", "Cup",
		"Plate", "Bowl", "Spoon", "Fork", "Knife",
		"Notebook", "Pen", "Pencil", "Marker", "Eraser",
		"Toy", "Game", "Puzzle", "Ball", "Doll",
		"Shirt", "Pants", "Shoes", "Hat", "Gloves",
	}
)

func generateRandomProduct() Product {
	productType := productTypes[rand.Intn(len(productTypes))]
	productName := productNames[rand.Intn(len(productNames))]
	category := categories[rand.Intn(len(categories))]

	name := fmt.Sprintf("%s %s", productType, productName)
	description := fmt.Sprintf("High-quality %s for your needs", name)
	price := float64(rand.Intn(99000)+1000) / 100.0 // $10.00 to $1000.00
	stock := rand.Intn(100) + 1
	imageURL := fmt.Sprintf("https://via.placeholder.com/300x300?text=%s", name)

	return Product{
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		Stock:       stock,
		ImageURL:    imageURL,
	}
}

func createProduct(catalogURL string, product Product) error {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	resp, err := http.Post(catalogURL+"/api/catalog/products", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	rand.Seed(time.Now().UnixNano())

	catalogURL := getEnv("CATALOG_SERVICE_URL", "http://localhost:8002")
	minProducts := 25
	maxProducts := 37

	loggerx.Info("Product generator started")
	loggerx.Infof("Target: %d-%d products per minute", minProducts, maxProducts)
	loggerx.Infof("Catalog service URL: %s", catalogURL)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Generate initial batch
	generateBatch(catalogURL, minProducts, maxProducts)

	for range ticker.C {
		generateBatch(catalogURL, minProducts, maxProducts)
	}
}

func generateBatch(catalogURL string, minProducts, maxProducts int) {
	productsToGenerate := rand.Intn(maxProducts-minProducts+1) + minProducts

	loggerx.Infof("Generating %d products...", productsToGenerate)

	successCount := 0
	failCount := 0

	for i := 0; i < productsToGenerate; i++ {
		product := generateRandomProduct()

		if err := createProduct(catalogURL, product); err != nil {
			loggerx.Errorf("Failed to create product: %v", err)
			failCount++
		} else {
			successCount++
		}

		// Small delay to avoid overwhelming the service
		time.Sleep(time.Duration(60000/productsToGenerate) * time.Millisecond)
	}

	loggerx.Infof("Batch complete: %d succeeded, %d failed", successCount, failCount)
}
