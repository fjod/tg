package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var db *sql.DB

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize database connection if not already done
	if db == nil {
		var err error
		db, err = initDB()
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}
	}
	// Get bot token from environment
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Printf("TELEGRAM_BOT_TOKEN not set")
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	// Create bot instance
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("Failed to create bot: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	// Parse incoming webhook
	var update tgbotapi.Update
	if err := json.Unmarshal([]byte(request.Body), &update); err != nil {
		log.Printf("Error parsing update: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	// Handle the message
	if update.Message != nil {
		handleMessage(bot, update.Message, db)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
