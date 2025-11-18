package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"mizon/loggerx"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type ValidateResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id,omitempty"`
}

func initDB() error {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "mizon")
	password := getEnv("POSTGRES_PASSWORD", "mizon123")
	dbname := getEnv("POSTGRES_DB", "mizon_users")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	for i := 0; i < 30; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				loggerx.Info("Successfully connected to PostgreSQL")
				return nil
			}
		}
		loggerx.Infof("Waiting for PostgreSQL... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user User
	err := db.QueryRow("SELECT id, username, email FROM users WHERE username = $1 AND password = $2",
		req.Username, req.Password).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session token
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = db.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)",
		user.ID, token, expiresAt)

	if err != nil {
		loggerx.Errorf("Error creating session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(ValidateResponse{Valid: false})
		return
	}

	var userID string
	var expiresAt time.Time
	err := db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = $1", token).
		Scan(&userID, &expiresAt)

	if err != nil || time.Now().After(expiresAt) {
		json.NewEncoder(w).Encode(ValidateResponse{Valid: false})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ValidateResponse{
		Valid:  true,
		UserID: userID,
	})
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
	if err := initDB(); err != nil {
		loggerx.Fatalf("%v", err)
	}
	defer db.Close()

	router := mux.NewRouter()
	cfg := loggerx.Config{LogRequestBody: loggerx.EnvBool("LOG_REQUEST_BODY", false), MaxBody: loggerx.EnvInt("LOG_MAX_BODY", 2048)}
	router.Use(loggerx.Middleware(cfg))
	router.HandleFunc("/api/auth/login", enableCORS(loginHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/validate", enableCORS(validateHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/health", enableCORS(healthHandler)).Methods("GET", "OPTIONS")

	port := getEnv("PORT", "8001")
	loggerx.Infof("Auth service starting on port %s", port)
	loggerx.Fatalf("%v", http.ListenAndServe(":"+port, router))
}
