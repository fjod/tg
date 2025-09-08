package main

import (
	"fmt"
	"log"
	"strings"

	telegramparser "github.com/kd3n1z/go-telegram-parser"
)

func validateTelegramWebApp(initData string, botToken string) (int64, error) {
	parser := telegramparser.CreateParser(botToken)
	validatedData, err := parser.Parse(initData)
	if err != nil {
		log.Printf("[WARN] Telegram WebApp validation failed: %v", err)
		return 0, fmt.Errorf("invalid initData: %v", err)
	}

	log.Printf("[INFO] Telegram WebApp validation successful")
	log.Printf("[INFO] User ID: %d, FirstName: %s", validatedData.User.Id, validatedData.User.FirstName)

	return validatedData.User.Id, nil
}

func extractUserIDFromAuth(authHeader string) (int64, error) {
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
