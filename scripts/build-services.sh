#!/bin/bash

echo "Building all Mizon services..."

# Build Auth Service
echo "Building Auth Service..."
cd services/auth
go mod download
go build -o auth-service .
cd ../..

# Build Catalog Service
echo "Building Catalog Service..."
cd services/catalog
go mod download
go build -o catalog-service .
cd ../..

# Build Search Service
echo "Building Search Service..."
cd services/search
go mod download
go build -o search-service .
cd ../..

# Build Cart Service
echo "Building Cart Service..."
cd services/cart
go mod download
go build -o cart-service .
cd ../..

# Build Order Service
echo "Building Order Service..."
cd services/order
go mod download
go build -o order-service .
cd ../..

echo "All services built successfully!"
