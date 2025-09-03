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
	log.Printf("Handler started - RequestID from context")

	// Initialize database connection if not already done
	if db == nil {
		log.Printf("Initializing database connection...")
		var err error
		db, err = initDB()
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}
		log.Printf("Database connection established")
	}

	// Get bot token from environment
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Printf("TELEGRAM_BOT_TOKEN not set")
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	// Create bot instance
	log.Printf("Creating bot instance...")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("Failed to create bot: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	// Parse incoming webhook
	log.Printf("Parsing webhook data...")
	var update tgbotapi.Update
	if err := json.Unmarshal([]byte(request.Body), &update); err != nil {
		log.Printf("Error parsing update: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	// Handle the message
	if update.Message != nil {
		log.Printf("Processing message from user %d", update.Message.From.ID)
		handleMessage(bot, update.Message, db)
	}

	// Handle callback queries (button clicks)
	if update.CallbackQuery != nil {
		log.Printf("Processing callback query from user %d", update.CallbackQuery.From.ID)
		handleCallbackQuery(bot, update.CallbackQuery, db)
	}

	log.Printf("Handler completed successfully")
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
