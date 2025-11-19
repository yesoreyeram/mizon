package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Request/Response types for new endpoints
type SignupRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type SignupResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type ProfileResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateProfileRequest struct {
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// Password hashing functions using bcrypt
func hashPassword(password string) (string, error) {
	// Use bcrypt cost 12 for good security/performance balance
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Password validation - enforce strong password requirements
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return errors.New("password must not exceed 128 characters")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// Email validation
func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	// Basic email regex - RFC 5322 simplified
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	if len(email) > 254 {
		return errors.New("email too long")
	}

	return nil
}

// Username validation
func validateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}

	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if len(username) > 50 {
		return errors.New("username must not exceed 50 characters")
	}

	// Allow alphanumeric, underscore, and hyphen
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}

// Input sanitization to prevent XSS
func sanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)
	// Escape HTML to prevent XSS
	input = html.EscapeString(input)
	return input
}

// Generate secure random token for password reset
func generateResetToken() (string, error) {
	// Generate 32 random bytes
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Encode to base64 URL-safe string
	return base64.URLEncoding.EncodeToString(b), nil
}

// Create session with optional remember me
func createSession(ctx context.Context, userID string, rememberMe bool) (string, time.Time, error) {
	if db == nil {
		return "", time.Time{}, errors.New("database not initialized")
	}

	token := uuid.New().String()
	var expiresAt time.Time

	if rememberMe {
		// 30 days for remember me
		expiresAt = time.Now().Add(30 * 24 * time.Hour)
	} else {
		// 24 hours for normal session
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	_, err := db.ExecContext(ctx,
		"INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, expiresAt)

	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to create session: %w", err)
	}

	return token, expiresAt, nil
}

// Get user ID from authorization token
func getUserIDFromToken(ctx context.Context, token string) (string, error) {
	if db == nil {
		return "", errors.New("database not initialized")
	}

	if token == "" {
		return "", errors.New("no token provided")
	}

	var userID string
	var expiresAt time.Time
	err := db.QueryRowContext(ctx,
		"SELECT user_id, expires_at FROM sessions WHERE token = $1",
		token).Scan(&userID, &expiresAt)

	if err != nil {
		return "", errors.New("invalid token")
	}

	if time.Now().After(expiresAt) {
		return "", errors.New("token expired")
	}

	return userID, nil
}

// Rate limiting helper (simple in-memory implementation)
type rateLimiter struct {
	attempts map[string][]time.Time
}

func newRateLimiter() *rateLimiter {
	return &rateLimiter{
		attempts: make(map[string][]time.Time),
	}
}

func (rl *rateLimiter) isAllowed(key string, maxAttempts int, window time.Duration) bool {
	now := time.Now()

	// Clean old attempts
	if attempts, exists := rl.attempts[key]; exists {
		var validAttempts []time.Time
		for _, t := range attempts {
			if now.Sub(t) < window {
				validAttempts = append(validAttempts, t)
			}
		}
		rl.attempts[key] = validAttempts
	}

	// Check if limit exceeded
	if len(rl.attempts[key]) >= maxAttempts {
		return false
	}

	// Add current attempt
	rl.attempts[key] = append(rl.attempts[key], now)
	return true
}

var loginRateLimiter = newRateLimiter()
var signupRateLimiter = newRateLimiter()
