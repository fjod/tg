package main

import (
	"fmt"
	"log"
	"strings"

	telegramparser "github.com/kd3n1z/go-telegram-parser"
)

type ParserInterface interface {
	Parse(query string) (telegramparser.WebAppInitData, error)
}

type ParserFactory func(botToken string) ParserInterface

var defaultParserFactory ParserFactory = func(botToken string) ParserInterface {
	parser := telegramparser.CreateParser(botToken)
	return &parser
}

func validateTelegramWebApp(initData string, p ParserInterface) (int64, error) {
	validatedData, err := p.Parse(initData)
	if err != nil {
		log.Printf("[WARN] Telegram WebApp validation failed: %v", err)
		return 0, fmt.Errorf("invalid initData: %v", err)
	}

	log.Printf("[INFO] Telegram WebApp validation successful")
	log.Printf("[INFO] User ID: %d, FirstName: %s", validatedData.User.Id, validatedData.User.FirstName)

	return validatedData.User.Id, nil
}

func extractUserIDFromAuth(authHeader string, envProvider EnvProvider, parserFactory ParserFactory) (int64, error) {
	if authHeader == "" {
		return 0, fmt.Errorf("authorization header is required")
	}

	// Remove "Bearer " prefix if present
	initData := strings.TrimPrefix(authHeader, "Bearer ")

	// Get bot token from environment
	botToken := envProvider.GetBotToken()
	if botToken == "" {
		return 0, fmt.Errorf("bot token not configured")
	}

	f := parserFactory(botToken)
	return validateTelegramWebApp(initData, f)
}
