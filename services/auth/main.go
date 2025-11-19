package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"mizon/loggerx"
	"mizon/telemetry"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me,omitempty"`
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
	driverName := "postgres"
	if dn, derr := telemetry.RegisterPostgres("postgres"); derr == nil {
		driverName = dn
	}
	for i := 0; i < 30; i++ {
		db, err = sql.Open(driverName, connStr)
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

	// Rate limiting - max 5 attempts per minute per IP
	clientIP := r.RemoteAddr
	if !loginRateLimiter.isAllowed(clientIP, 5, time.Minute) {
		http.Error(w, "Too many login attempts. Please try again later.", http.StatusTooManyRequests)
		return
	}

	// Sanitize inputs
	username := sanitizeInput(req.Username)

	// Query user with password hash
	var user User
	var passwordHash string
	err := db.QueryRowContext(r.Context(),
		"SELECT id, username, email, password FROM users WHERE username = $1",
		username).Scan(&user.ID, &user.Username, &user.Email, &passwordHash)

	if err != nil {
		// Use constant time comparison to prevent timing attacks
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password using bcrypt
	if !checkPasswordHash(req.Password, passwordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session with remember me support
	token, _, err := createSession(r.Context(), user.ID, req.RememberMe)
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
	err := db.QueryRowContext(r.Context(), "SELECT user_id, expires_at FROM sessions WHERE token = $1", token).
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

// signupHandler handles user registration
func signupHandler(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Rate limiting - max 3 signups per hour per IP
	clientIP := r.RemoteAddr
	if !signupRateLimiter.isAllowed(clientIP, 3, time.Hour) {
		http.Error(w, "Too many signup attempts. Please try again later.", http.StatusTooManyRequests)
		return
	}

	// Sanitize inputs
	req.Username = sanitizeInput(req.Username)
	req.Email = sanitizeInput(req.Email)
	req.FirstName = sanitizeInput(req.FirstName)
	req.LastName = sanitizeInput(req.LastName)

	// Validate username
	if err := validateUsername(req.Username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate email
	if err := validateEmail(req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate password strength
	if err := validatePassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		loggerx.Errorf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if username already exists
	var existingID string
	err = db.QueryRowContext(r.Context(), "SELECT id FROM users WHERE username = $1", req.Username).Scan(&existingID)
	if err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Check if email already exists
	err = db.QueryRowContext(r.Context(), "SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
	if err == nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Insert new user
	var userID string
	err = db.QueryRowContext(r.Context(),
		`INSERT INTO users (username, password, email, first_name, last_name) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Username, passwordHash, req.Email, req.FirstName, req.LastName).Scan(&userID)

	if err != nil {
		loggerx.Errorf("Error creating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := SignupResponse{
		UserID:   userID,
		Username: req.Username,
		Email:    req.Email,
		Message:  "User created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// forgotPasswordHandler initiates password reset process
func forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.Email = sanitizeInput(req.Email)

	// Validate email format
	if err := validateEmail(req.Email); err != nil {
		// Always return success to prevent email enumeration
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ForgotPasswordResponse{
			Message: "If the email exists, a password reset link has been sent",
		})
		return
	}

	// Check if user exists
	var userID string
	err := db.QueryRowContext(r.Context(), "SELECT id FROM users WHERE email = $1", req.Email).Scan(&userID)
	if err != nil {
		// Don't reveal if email doesn't exist - always return success
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ForgotPasswordResponse{
			Message: "If the email exists, a password reset link has been sent",
		})
		return
	}

	// Generate reset token
	resetToken, err := generateResetToken()
	if err != nil {
		loggerx.Errorf("Error generating reset token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store reset token (expires in 1 hour)
	expiresAt := time.Now().Add(time.Hour)
	_, err = db.ExecContext(r.Context(),
		`INSERT INTO password_reset_tokens (user_id, token, expires_at) 
		 VALUES ($1, $2, $3)`,
		userID, resetToken, expiresAt)

	if err != nil {
		loggerx.Errorf("Error storing reset token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// In production, send email with reset link here
	loggerx.Infof("Password reset token for user %s: %s", userID, resetToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ForgotPasswordResponse{
		Message: "If the email exists, a password reset link has been sent",
	})
}

// resetPasswordHandler completes password reset
func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate new password
	if err := validatePassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify reset token
	var userID string
	var expiresAt time.Time
	err := db.QueryRowContext(r.Context(),
		"SELECT user_id, expires_at FROM password_reset_tokens WHERE token = $1",
		req.Token).Scan(&userID, &expiresAt)

	if err != nil {
		http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
		return
	}

	if time.Now().After(expiresAt) {
		http.Error(w, "Reset token has expired", http.StatusBadRequest)
		return
	}

	// Hash new password
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		loggerx.Errorf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update password
	_, err = db.ExecContext(r.Context(),
		"UPDATE users SET password = $1 WHERE id = $2",
		passwordHash, userID)

	if err != nil {
		loggerx.Errorf("Error updating password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Delete used token
	_, err = db.ExecContext(r.Context(),
		"DELETE FROM password_reset_tokens WHERE token = $1",
		req.Token)

	if err != nil {
		loggerx.Warnf("Error deleting reset token: %v", err)
	}

	// Invalidate all existing sessions for security
	_, err = db.ExecContext(r.Context(),
		"DELETE FROM sessions WHERE user_id = $1",
		userID)

	if err != nil {
		loggerx.Warnf("Error invalidating sessions: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResetPasswordResponse{
		Message: "Password reset successful",
	})
}

// profileHandler returns user profile information
func profileHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user ID from token
	userID, err := getUserIDFromToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch user profile
	var profile ProfileResponse
	err = db.QueryRowContext(r.Context(),
		`SELECT id, username, email, first_name, last_name, created_at 
		 FROM users WHERE id = $1`,
		userID).Scan(&profile.ID, &profile.Username, &profile.Email,
		&profile.FirstName, &profile.LastName, &profile.CreatedAt)

	if err != nil {
		loggerx.Errorf("Error fetching profile: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// updateProfileHandler updates user profile information
func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user ID from token
	userID, err := getUserIDFromToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Sanitize inputs
	req.Email = sanitizeInput(req.Email)
	req.FirstName = sanitizeInput(req.FirstName)
	req.LastName = sanitizeInput(req.LastName)

	// Validate email if provided
	if req.Email != "" {
		if err := validateEmail(req.Email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if email already exists for another user
		var existingID string
		err = db.QueryRowContext(r.Context(),
			"SELECT id FROM users WHERE email = $1 AND id != $2",
			req.Email, userID).Scan(&existingID)
		if err == nil {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}
	}

	// Build dynamic update query
	query := "UPDATE users SET "
	args := []interface{}{}
	argCount := 1

	if req.Email != "" {
		query += fmt.Sprintf("email = $%d, ", argCount)
		args = append(args, req.Email)
		argCount++
	}
	if req.FirstName != "" {
		query += fmt.Sprintf("first_name = $%d, ", argCount)
		args = append(args, req.FirstName)
		argCount++
	}
	if req.LastName != "" {
		query += fmt.Sprintf("last_name = $%d, ", argCount)
		args = append(args, req.LastName)
		argCount++
	}

	// Remove trailing comma and space
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, userID)

	// Execute update
	_, err = db.ExecContext(r.Context(), query, args...)
	if err != nil {
		loggerx.Errorf("Error updating profile: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile updated successfully",
	})
}

// logoutHandler invalidates the current session
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete the session
	_, err := db.ExecContext(r.Context(), "DELETE FROM sessions WHERE token = $1", token)
	if err != nil {
		loggerx.Errorf("Error deleting session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
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
	if _, err := telemetry.Setup("auth-service"); err != nil {
		loggerx.Warnf("tracing setup failed: %v", err)
	}
	if err := initDB(); err != nil {
		loggerx.Fatalf("%v", err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.Use(telemetry.MuxMiddleware("auth-service"))
	cfg := loggerx.Config{LogRequestBody: loggerx.EnvBool("LOG_REQUEST_BODY", false), MaxBody: loggerx.EnvInt("LOG_MAX_BODY", 2048)}
	router.Use(loggerx.Middleware(cfg))

	// Metrics endpoint
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	// Health check
	router.HandleFunc("/health", enableCORS(healthHandler)).Methods("GET", "OPTIONS")

	// Authentication endpoints
	router.HandleFunc("/api/auth/signup", enableCORS(signupHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", enableCORS(loginHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/logout", enableCORS(logoutHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/validate", enableCORS(validateHandler)).Methods("GET", "OPTIONS")

	// Password reset endpoints
	router.HandleFunc("/api/auth/forgot-password", enableCORS(forgotPasswordHandler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/reset-password", enableCORS(resetPasswordHandler)).Methods("POST", "OPTIONS")

	// Profile endpoints
	router.HandleFunc("/api/auth/profile", enableCORS(profileHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/auth/profile", enableCORS(updateProfileHandler)).Methods("PUT", "OPTIONS")

	port := getEnv("PORT", "8001")
	loggerx.Infof("Auth service starting on port %s", port)
	loggerx.Fatalf("%v", http.ListenAndServe(":"+port, router))
}
