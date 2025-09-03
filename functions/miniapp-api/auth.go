package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
}

func validateTelegramWebApp(initData string, botToken string) (int64, error) {
	// Parse the initData parameters
	values, err := url.ParseQuery(initData)
	if err != nil {
		return 0, fmt.Errorf("failed to parse initData: %v", err)
	}

	// Extract the hash
	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return 0, fmt.Errorf("hash parameter is missing")
	}

	// Remove hash from values for validation
	values.Del("hash")

	// Create data check string
	var pairs []string
	for key, valueSlice := range values {
		if len(valueSlice) > 0 {
			pairs = append(pairs, key+"="+valueSlice[0])
		}
	}
	sort.Strings(pairs)
	dataCheckString := strings.Join(pairs, "\n")

	// Create secret key using bot token
	secretKey := sha256.Sum256([]byte(botToken))

	// Calculate expected hash
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	// Verify hash
	if receivedHash != expectedHash {
		return 0, fmt.Errorf("invalid hash")
	}

	// Extract user information
	userStr := values.Get("user")
	if userStr == "" {
		return 0, fmt.Errorf("user parameter is missing")
	}

	var user TelegramUser
	if err := json.Unmarshal([]byte(userStr), &user); err != nil {
		return 0, fmt.Errorf("failed to parse user data: %v", err)
	}

	return user.ID, nil
}

func extractUserIDFromAuth(authHeader string) (int64, error) {
	// authHeader should contain the initData from Telegram Web App
	if authHeader == "" {
		return 0, fmt.Errorf("authorization header is required")
	}

	// Remove "Bearer " prefix if present
	initData := strings.TrimPrefix(authHeader, "Bearer ")

	// Get bot token from environment
	botToken := getBotToken()
	if botToken == "" {
		return 0, fmt.Errorf("bot token not configured")
	}

	return validateTelegramWebApp(initData, botToken)
}
