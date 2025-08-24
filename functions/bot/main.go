package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
		handleMessage(bot, update.Message)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	var responseText string

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			responseText = "Hello! I'm your Telegram Content Organizer bot. Send me any message or forward content to me!"
		case "help":
			responseText = "Available commands:\n/start - Get started\n/help - Show this help message\n\nYou can also send me any message or forward content to me."
		default:
			responseText = "Unknown command. Use /help to see available commands."
		}
	} else if message.ForwardFrom != nil {
		responseText = "Thanks for forwarding this message! I've received it and will organize it for you."
		log.Printf("Forwarded message from %s: %s", message.ForwardFrom.FirstName, message.Text)
	} else {
		responseText = "Hello Nigger!!! You said: " + message.Text
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ReplyToMessageID = message.MessageID

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func main() {
	lambda.Start(Handler)
}
