package main

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	// Save user to database
	if err := saveUser(db, message.From); err != nil {
		log.Printf("Error saving user: %v", err)
	}

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
	} else {
		// Save message to database for all non-command messages
		if err := saveMessage(db, message); err != nil {
			log.Printf("Error saving message: %v", err)
			responseText = "Sorry, I couldn't save your message. Please try again."
		} else {
			if message.ForwardFrom != nil {
				responseText = "Thanks for forwarding this message! I've saved it and will organize it for you."
				log.Printf("Forwarded message from %s saved successfully", message.ForwardFrom.FirstName)
			} else {
				responseText = "Message saved! I've organized it for you."
			}
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ReplyToMessageID = message.MessageID

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
