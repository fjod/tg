package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Tag struct {
	ID        int64     `json:"id"         db:"id"`
	UserID    int64     `json:"user_id"    db:"user_id"`
	Name      string    `json:"name"       db:"name"`
	Color     *string   `json:"color"      db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func getUserTags(db *sql.DB, userID int64) ([]Tag, error) {
	query := `SELECT id, name, color FROM tags WHERE user_id = $1 ORDER BY name`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		var color sql.NullString
		if err := rows.Scan(&tag.ID, &tag.Name, &color); err != nil {
			return nil, err
		}
		tag.UserID = userID
		if color.Valid {
			tag.Color = &color.String
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func getOrCreateTag(db *sql.DB, userID int64, tagName string) (int64, error) {
	var tagID int64

	// Try to get existing tag
	query := `SELECT id FROM tags WHERE user_id = $1 AND name = $2`
	err := db.QueryRow(query, userID, tagName).Scan(&tagID)

	if err == sql.ErrNoRows {
		// Create new tag
		insertQuery := `INSERT INTO tags (user_id, name, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP) RETURNING id`
		err = db.QueryRow(insertQuery, userID, tagName).Scan(&tagID)
	}

	return tagID, err
}

func tagMessage(db *sql.DB, messageID int64, tagID int64) error {
	query := `INSERT INTO message_tags (message_id, tag_id, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP) ON CONFLICT (message_id, tag_id) DO NOTHING`
	_, err := db.Exec(query, messageID, tagID)
	return err
}

func getMessageByTelegramID(db *sql.DB, userID int64, telegramMessageID int64) (int64, error) {
	var messageID int64
	query := `SELECT id FROM messages WHERE user_id = $1 AND telegram_message_id = $2`
	err := db.QueryRow(query, userID, telegramMessageID).Scan(&messageID)
	return messageID, err
}

func showTagSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	// Get user's existing tags
	tags, err := getUserTags(db, message.From.ID)
	if err != nil {
		log.Printf("Error getting user tags: %v", err)
		sendErrorMessage(bot, message, "Could not load your tags.")
		return
	}

	var responseText string
	if len(tags) == 0 {
		responseText = fmt.Sprintf("You don't have any tags yet. Reply with a tag name to create your first tag:\n\n[MSG_ID:%d]", message.MessageID)
	} else {
		responseText = fmt.Sprintf("Choose a tag by typing its name or create a new one:\n\n")
		for i, tag := range tags {
			responseText += fmt.Sprintf("%d. %s\n", i+1, tag.Name)
		}
		responseText += fmt.Sprintf("\nOr type a new tag name to create it.\n\n[MSG_ID:%d]", message.MessageID)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending tag selection: %v", err)
	}
}

func handleTagSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	// Extract original message ID from the bot's tag selection message
	if message.ReplyToMessage == nil {
		log.Printf("No ReplyToMessage found")
		sendErrorMessage(bot, message, "This doesn't appear to be a reply.")
		return
	}
	
	// Parse the original message ID from the tag selection message text
	botMessageText := message.ReplyToMessage.Text
	msgIDStart := strings.Index(botMessageText, "[MSG_ID:")
	if msgIDStart == -1 {
		log.Printf("Could not find MSG_ID in bot message: %s", botMessageText)
		sendErrorMessage(bot, message, "Could not find the original message to tag.")
		return
	}
	
	msgIDEnd := strings.Index(botMessageText[msgIDStart:], "]")
	if msgIDEnd == -1 {
		log.Printf("Could not find closing bracket for MSG_ID")
		sendErrorMessage(bot, message, "Could not find the original message to tag.")
		return
	}
	
	msgIDStr := botMessageText[msgIDStart+9 : msgIDStart+msgIDEnd] // +9 to skip "[MSG_ID:"
	originalMessageID, err := strconv.Atoi(msgIDStr)
	if err != nil {
		log.Printf("Could not parse message ID: %s", msgIDStr)
		sendErrorMessage(bot, message, "Could not find the original message to tag.")
		return
	}
	
	log.Printf("Extracted original message ID: %d", originalMessageID)

	// Get the database message ID
	dbMessageID, err := getMessageByTelegramID(db, message.From.ID, int64(originalMessageID))
	if err != nil {
		log.Printf("Error finding original message: %v", err)
		sendErrorMessage(bot, message, "Could not find the original message to tag.")
		return
	}

	// Parse tag selection
	tagName := strings.TrimSpace(message.Text)
	if tagName == "" {
		sendErrorMessage(bot, message, "Please enter a tag name.")
		return
	}

	// Check if it's a number (selecting from list)
	if num, err := strconv.Atoi(tagName); err == nil {
		// User selected by number
		tags, err := getUserTags(db, message.From.ID)
		if err != nil || num < 1 || num > len(tags) {
			sendErrorMessage(bot, message, "Invalid tag number. Please try again.")
			return
		}
		tagName = tags[num-1].Name
	}

	// Get or create the tag
	tagID, err := getOrCreateTag(db, message.From.ID, tagName)
	if err != nil {
		log.Printf("Error creating/getting tag: %v", err)
		sendErrorMessage(bot, message, "Could not create or find the tag.")
		return
	}

	// Tag the message
	if err := tagMessage(db, dbMessageID, tagID); err != nil {
		log.Printf("Error tagging message: %v", err)
		sendErrorMessage(bot, message, "Could not tag the message.")
		return
	}

	// Send confirmation
	responseText := fmt.Sprintf("✅ Message tagged with '%s'", tagName)
	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending confirmation: %v", err)
	}
}

func sendErrorMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending error message: %v", err)
	}
}
