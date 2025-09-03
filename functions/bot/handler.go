package main

import (
	"database/sql"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	// Save user to database
	if err := saveUser(db, message.From); err != nil {
		log.Printf("Error saving user : %v", err)
	}

	var responseText string

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			responseText = "Hello! I'm your Telegram Content Organizer bot. Send me any message or forward content to me!"
		case "help":
			responseText = "Available commands:\n/start - Get started\n/help - Show this help message\n/miniapp - Open mini-app to view your tags\n\nYou can also send me any message or forward content to me."
		case "miniapp":
			sendMiniAppButton(bot, message)
			return
		default:
			responseText = "Unknown command. Use /help to see available commands."
		}
	} else {
		// Check if this is a reply to our tag selection message
		if message.ReplyToMessage != nil && message.ReplyToMessage.From.IsBot {
			// Check if the reply is to a tag selection message by checking message content
			if strings.Contains(message.ReplyToMessage.Text, "Choose a tag by typing") ||
				strings.Contains(message.ReplyToMessage.Text, "You don't have any tags yet") ||
				strings.Contains(message.ReplyToMessage.Text, "[MSG_ID:") {
				handleTagSelection(bot, message, db)
				return
			}
		}

		// Save message to database for all non-command messages
		if err := saveMessage(db, message); err != nil {
			log.Printf("Error saving message: %v", err)
			responseText = "Sorry, I couldn't save your message. Please try again."
		} else {
			// Show tag selection after saving message
			showTagSelection(bot, message, db)
			return
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ReplyToMessageID = message.MessageID

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, db *sql.DB) {
	// Answer the callback query to stop the loading animation
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := bot.Request(callback); err != nil {
		log.Printf("Error answering callback query: %v", err)
	}

	// Parse callback data format: "tag:tagID:messageID" or "new_tag:messageID"
	data := callbackQuery.Data
	log.Printf("Received callback data: %s", data)

	if strings.HasPrefix(data, "tag:") {
		handleTagCallback(bot, callbackQuery, db)
	} else if strings.HasPrefix(data, "new_tag:") {
		handleNewTagCallback(bot, callbackQuery, db)
	} else {
		log.Printf("Unknown callback data format: %s", data)
	}
}

func sendMiniAppButton(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Create a keyboard with URL button (opens mini-app)
	webAppURL := "https://functions.yandexcloud.net/d4ek5oug8uak4lb9edsl"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🏷️ View My Tags", webAppURL),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Open the mini-app to view and manage your tags:")
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending mini-app button: %v", err)
	}
}
