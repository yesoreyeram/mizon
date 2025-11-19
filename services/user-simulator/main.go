package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"mizon/loggerx"
)

type SignupRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

var (
	firstNames = []string{
		"Emma", "Liam", "Olivia", "Noah", "Ava", "Ethan", "Sophia", "Mason",
		"Isabella", "William", "Mia", "James", "Charlotte", "Benjamin", "Amelia",
		"Lucas", "Harper", "Henry", "Evelyn", "Alexander", "Abigail", "Michael",
		"Emily", "Daniel", "Elizabeth", "Matthew", "Sofia", "Joseph", "Avery",
		"David", "Ella", "Jackson", "Madison", "Logan", "Scarlett", "Sebastian",
		"Victoria", "Jack", "Aria", "Aiden", "Grace", "Owen", "Chloe", "Samuel",
		"Camila", "Gabriel", "Penelope", "Carter", "Riley", "Jayden",
	}

	lastNames = []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller",
		"Davis", "Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez",
		"Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin",
		"Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark",
		"Ramirez", "Lewis", "Robinson", "Walker", "Young", "Allen", "King",
		"Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores", "Green",
		"Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell",
		"Carter", "Roberts",
	}

	passwordChars = []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}

	specialChars = []string{"!", "@", "#", "$", "%", "^", "&", "*"}
)

func generateRandomUser() SignupRequest {
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	timestamp := time.Now().UnixNano()

	// Generate unique username with timestamp
	username := fmt.Sprintf("%s%s%d", 
		firstName[:min(4, len(firstName))], 
		lastName[:min(4, len(lastName))], 
		timestamp%100000)

	// Generate email
	email := fmt.Sprintf("%s.%s%d@example.com", 
		firstName, 
		lastName, 
		timestamp%10000)

	// Generate secure password that meets requirements
	password := generateSecurePassword()

	return SignupRequest{
		Username:  username,
		Password:  password,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
}

func generateSecurePassword() string {
	// Generate a password that meets requirements:
	// - At least 8 characters
	// - Contains uppercase
	// - Contains lowercase
	// - Contains number
	// - Contains special character
	
	password := ""
	
	// Add at least one of each required type
	password += string(passwordChars[rand.Intn(26)]) // Uppercase (A-Z)
	password += string(passwordChars[26+rand.Intn(26)]) // Lowercase (a-z)
	password += string(passwordChars[52+rand.Intn(10)]) // Number (0-9)
	password += specialChars[rand.Intn(len(specialChars))] // Special char
	
	// Fill the rest to make at least 12 characters
	for len(password) < 12 {
		if rand.Intn(10) < 2 {
			password += specialChars[rand.Intn(len(specialChars))]
		} else {
			password += passwordChars[rand.Intn(len(passwordChars))]
		}
	}
	
	// Shuffle the password
	runes := []rune(password)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	
	return string(runes)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func signupUser(authURL string, user SignupRequest) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	resp, err := http.Post(authURL+"/api/auth/signup", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("rate limit exceeded (429)")
	}

	if resp.StatusCode == http.StatusConflict {
		return fmt.Errorf("user already exists (409)")
	}

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

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func main() {
	rand.Seed(time.Now().UnixNano())

	authURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8001")
	minUsers := getEnvInt("MIN_USERS_PER_MINUTE", 2)
	maxUsers := getEnvInt("MAX_USERS_PER_MINUTE", 7)

	// Validate constraints
	if minUsers < 1 {
		minUsers = 1
	}
	if maxUsers < minUsers {
		maxUsers = minUsers
	}
	if minUsers > 10 || maxUsers > 10 {
		loggerx.Warn("High user generation rate may trigger rate limits")
	}

	loggerx.Info("User simulator started")
	loggerx.Infof("Target: %d-%d users per minute", minUsers, maxUsers)
	loggerx.Infof("Auth service URL: %s", authURL)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Generate initial batch
	generateBatch(authURL, minUsers, maxUsers)

	for range ticker.C {
		generateBatch(authURL, minUsers, maxUsers)
	}
}

func generateBatch(authURL string, minUsers, maxUsers int) {
	usersToGenerate := rand.Intn(maxUsers-minUsers+1) + minUsers

	loggerx.Infof("Generating %d users...", usersToGenerate)

	successCount := 0
	failCount := 0
	rateLimitCount := 0

	for i := 0; i < usersToGenerate; i++ {
		user := generateRandomUser()

		if err := signupUser(authURL, user); err != nil {
			if err.Error() == "rate limit exceeded (429)" {
				rateLimitCount++
			} else {
				loggerx.Errorf("Failed to create user: %v", err)
			}
			failCount++
		} else {
			successCount++
			loggerx.Infof("Created user: %s (%s)", user.Username, user.Email)
		}

		// Spread requests across the minute to avoid bursts
		time.Sleep(time.Duration(60000/usersToGenerate) * time.Millisecond)
	}

	loggerx.Infof("Batch complete: %d succeeded, %d failed", successCount, failCount)
	if rateLimitCount > 0 {
		loggerx.Warnf("Rate limit hit %d times - consider reducing MAX_USERS_PER_MINUTE", rateLimitCount)
	}
}
