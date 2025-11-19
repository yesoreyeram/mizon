# Mizon - Minimal E-commerce Platform

A microservices-based e-commerce platform similar to Amazon with minimal functionality for demonstration purposes.

## Architecture Overview

This application demonstrates a microservices architecture with the following services:

### Backend Services (Go 1.23)

- **Auth Service** (Port 8001): User authentication using PostgreSQL
- **Catalog Service** (Port 8002): Product catalog management using MongoDB
- **Search Service** (Port 8003): Product search using Elasticsearch
- **Cart Service** (Port 8004): Shopping cart management using MongoDB
- **Order Service** (Port 8005): Order processing using Kafka messaging

### Frontend

- **Next.js App** (Port 3000): React-based UI with TypeScript and TailwindCSS

### Infrastructure

- **PostgreSQL** (Port 5432): User data storage
- **MongoDB** (Port 27017): Product catalog and cart storage
- **Elasticsearch** (Port 9200): Product search indexing
- **Kafka** (Port 9092): Message queue for order processing
- **Zookeeper** (Port 2181): Kafka coordination

## Technology Stack

- **Backend**: Go 1.23
- **Frontend**: Next.js, React, TypeScript, TailwindCSS
- **Databases**: PostgreSQL, MongoDB
- **Search**: Elasticsearch
- **Messaging**: Apache Kafka
- **Containerization**: Docker & Docker Compose

## Features

- ✅ **Secure User Authentication**
  - User registration with strong password requirements
  - Signin with remember me functionality
  - Password reset via email (token-based)
  - User profile management
  - Bcrypt password hashing
  - Rate limiting on auth endpoints
  - Session management
- ✅ Product Catalog Browsing
- ✅ Category-based Navigation
- ✅ Product Search
- ✅ Shopping Cart Management
- ✅ Order Placement
- ✅ 100 Pre-seeded Products

## Prerequisites

- Docker Desktop (with Docker Compose)
- Go 1.23+ (for local development)
- Node.js 18+ (for local development)

## Quick Start

### 1. Start Infrastructure Services

```bash
docker-compose up -d postgres mongodb elasticsearch zookeeper kafka
```

Wait for all services to be healthy (about 30 seconds):

```bash
docker-compose ps
```

### 2. Build and Start Backend Services

```bash
# Build all services
./scripts/build-services.sh

# Start all backend services
docker-compose up -d auth-service catalog-service search-service cart-service order-service
```

### 3. Bootstrap Catalog Data

```bash
# Run the bootstrap script to seed 100 products
./scripts/bootstrap.sh
```

### 4. Start Frontend

```bash
docker-compose up -d frontend
```

### 5. Access the Application

Open your browser to: http://localhost:3000

**Default Admin Credentials:**

- Username: `admin`
- Password: `admin`

## Alternative: Start Everything at Once

```bash
# Start all services
docker-compose up -d

# Wait for services to be ready (60 seconds)
sleep 60

# Bootstrap catalog
./scripts/bootstrap.sh
```

## API Endpoints

### Auth Service (http://localhost:8001)

- `POST /api/auth/signup` - Register a new user
- `POST /api/auth/login` - Login with username/password (supports remember me)
- `POST /api/auth/logout` - Logout and invalidate session
- `GET /api/auth/validate` - Validate token
- `GET /api/auth/profile` - Get user profile
- `PUT /api/auth/profile` - Update user profile
- `POST /api/auth/forgot-password` - Request password reset
- `POST /api/auth/reset-password` - Reset password with token

📖 **[Complete Auth API Documentation](docs/AUTH_API.md)**

### Catalog Service (http://localhost:8002)

- `GET /api/catalog/products` - List all products
- `GET /api/catalog/products/:id` - Get product details
- `GET /api/catalog/categories` - List categories
- `GET /api/catalog/categories/:category/products` - Products by category
- `POST /api/catalog/products` - Create product (admin only)

### Search Service (http://localhost:8003)

- `GET /api/search?q=keyword` - Search products
- `POST /api/search/index` - Index a product

### Cart Service (http://localhost:8004)

- `GET /api/cart/:userId` - Get user's cart
- `POST /api/cart/:userId/items` - Add item to cart
- `PUT /api/cart/:userId/items/:itemId` - Update item quantity
- `DELETE /api/cart/:userId/items/:itemId` - Remove item from cart
- `DELETE /api/cart/:userId` - Clear cart

### Order Service (http://localhost:8005)

- `POST /api/orders` - Create order (publishes to Kafka)
- `GET /api/orders/:userId` - Get user's orders
- `GET /api/orders/details/:orderId` - Get order details

## Project Structure

```txt
mizon/
├── services/
│   ├── auth/              # Authentication service
│   ├── catalog/           # Product catalog service
│   ├── search/            # Search service
│   ├── cart/              # Shopping cart service
│   └── order/             # Order processing service
├── frontend/              # Next.js frontend application
├── scripts/
│   ├── build-services.sh  # Build all Go services
│   ├── bootstrap.sh       # Seed initial data
│   └── init-db.sql        # PostgreSQL initialization
├── docker-compose.yml     # Docker orchestration
└── README.md             # This file
```

## Development

### Running Services Locally (without Docker)

#### Backend Services

```bash
# Terminal 1 - Auth Service
cd services/auth
go run main.go

# Terminal 2 - Catalog Service
cd services/catalog
go run main.go

# Terminal 3 - Search Service
cd services/search
go run main.go

# Terminal 4 - Cart Service
cd services/cart
go run main.go

# Terminal 5 - Order Service
cd services/order
go run main.go
```

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Environment Variables

Each service uses the following default environment variables (configured in docker-compose.yml):

**Auth Service:**

- `PORT=8001`
- `POSTGRES_HOST=postgres`
- `POSTGRES_PORT=5432`
- `POSTGRES_USER=mizon`
- `POSTGRES_PASSWORD=mizon123`
- `POSTGRES_DB=mizon_users`

**Catalog Service:**

- `PORT=8002`
- `MONGO_URI=mongodb://mongodb:27017`
- `MONGO_DB=mizon_catalog`

**Search Service:**

- `PORT=8003`
- `ELASTICSEARCH_URL=http://elasticsearch:9200`

**Cart Service:**

- `PORT=8004`
- `MONGO_URI=mongodb://mongodb:27017`
- `MONGO_DB=mizon_cart`

**Order Service:**

- `PORT=8005`
- `KAFKA_BROKER=kafka:9092`
- `MONGO_URI=mongodb://mongodb:27017`
- `MONGO_DB=mizon_orders`

**Frontend:**

- `NEXT_PUBLIC_AUTH_API=http://localhost:8001`
- `NEXT_PUBLIC_CATALOG_API=http://localhost:8002`
- `NEXT_PUBLIC_SEARCH_API=http://localhost:8003`
- `NEXT_PUBLIC_CART_API=http://localhost:8004`
- `NEXT_PUBLIC_ORDER_API=http://localhost:8005`

## Data Models

### User (PostgreSQL)

```json
{
  "id": "uuid",
  "username": "string",
  "password": "string (hashed)",
  "email": "string",
  "created_at": "timestamp"
}
```

### Product (MongoDB)

```json
{
  "_id": "ObjectId",
  "name": "string",
  "description": "string",
  "price": "number",
  "category": "string",
  "stock": "number",
  "image_url": "string",
  "created_at": "timestamp"
}
```

### Cart (MongoDB)

```json
{
  "_id": "ObjectId",
  "user_id": "string",
  "items": [
    {
      "product_id": "string",
      "name": "string",
      "price": "number",
      "quantity": "number"
    }
  ],
  "updated_at": "timestamp"
}
```

### Order (MongoDB + Kafka)

```json
{
  "_id": "ObjectId",
  "user_id": "string",
  "items": "array",
  "total": "number",
  "status": "string",
  "created_at": "timestamp"
}
```

## Troubleshooting

### Services Not Starting

```bash
# Check service logs
docker-compose logs -f [service-name]

# Restart specific service
docker-compose restart [service-name]
```

### Elasticsearch Not Ready

Elasticsearch takes about 30 seconds to start. Check status:

```bash
curl http://localhost:9200/_cluster/health
```

### Kafka Connection Issues

Ensure Zookeeper is running before Kafka:

```bash
docker-compose up -d zookeeper
sleep 10
docker-compose up -d kafka
```

### Database Connection Errors

Verify databases are accessible:

```bash
# PostgreSQL
docker-compose exec postgres psql -U mizon -d mizon_users -c "SELECT 1;"

# MongoDB
docker-compose exec mongodb mongosh --eval "db.adminCommand('ping')"
```

### Reset Everything

```bash
# Stop and remove all containers, volumes, and networks
docker-compose down -v

# Start fresh
docker-compose up -d
./scripts/bootstrap.sh
```

## Testing the Application

### 1. Test Authentication

```bash
curl -X POST http://localhost:8001/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'
```

### 2. Browse Products

```bash
curl http://localhost:8002/api/catalog/products
```

### 3. Search Products

```bash
curl "http://localhost:8003/api/search?q=laptop"
```

### 4. Add to Cart

```bash
curl -X POST http://localhost:8004/api/cart/admin/items \
  -H "Content-Type: application/json" \
  -d '{"product_id":"...","quantity":1}'
```

### 5. Place Order

```bash
curl -X POST http://localhost:8005/api/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":"admin","items":[...],"total":99.99}'
```

## Scaling

Individual services can be scaled using Docker Compose:

```bash
docker-compose up -d --scale catalog-service=3
```

## Stopping the Application

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (clears all data)
docker-compose down -v
```

## Future Enhancements

- ~~JWT-based authentication~~ ✅ (Token-based auth implemented)
- Payment gateway integration
- ~~User registration and profile management~~ ✅ (Implemented)
- Email verification for new accounts
- Two-factor authentication (2FA)
- OAuth integration (Google, GitHub, etc.)
- Product reviews and ratings
- Recommendation engine
- Admin dashboard
- Order tracking
- Email notifications
- Redis caching layer
- API Gateway (Kong/Nginx)
- Kubernetes deployment
- Monitoring (Prometheus/Grafana)
- CI/CD pipeline

## License

MIT License - For demonstration purposes only.

## Support

For issues and questions, please check the troubleshooting section above or review service logs using `docker-compose logs`.
