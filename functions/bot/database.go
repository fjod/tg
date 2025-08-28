package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}

func saveUser(db *sql.DB, user *tgbotapi.User) error {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, true)
		ON CONFLICT (telegram_id) 
		DO UPDATE SET 
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			updated_at = CURRENT_TIMESTAMP`

	var username, firstName, lastName sql.NullString
	if user.UserName != "" {
		username = sql.NullString{String: user.UserName, Valid: true}
	}
	if user.FirstName != "" {
		firstName = sql.NullString{String: user.FirstName, Valid: true}
	}
	if user.LastName != "" {
		lastName = sql.NullString{String: user.LastName, Valid: true}
	}

	_, err := db.Exec(query, user.ID, username, firstName, lastName)
	return err
}

func saveMessage(db *sql.DB, message *tgbotapi.Message) error {

	var textContent, caption sql.NullString

	// Store only previews/snippets of text content
	if message.Text != "" {
		preview := truncateText(message.Text, 150)
		textContent = sql.NullString{String: preview, Valid: true}
	}
	if message.Caption != "" {
		preview := truncateText(message.Caption, 150)
		caption = sql.NullString{String: preview, Valid: true}
	}

	// Extract file metadata
	messageType := getMessageType(message)
	fileMetadata := extractFileMetadata(message, messageType)

	// Extract metadata from FULL text and caption (not just previews)
	urls := extractURLs(message.Text, message.Caption)
	hashtags := extractHashtags(message.Text, message.Caption)
	mentions := extractMentions(message.Text, message.Caption)

	// Handle forwarded message data
	forwardedDate, forwardedFrom := generateForwardedTimes(message)

	query := `
		INSERT INTO messages (
			user_id, telegram_message_id, message_type, text_content, caption,
			file_id, file_name, file_size, mime_type, duration,
			forwarded_date, forwarded_from, urls, hashtags, mentions, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, CURRENT_TIMESTAMP)`

	_, err := db.Exec(query,
		message.From.ID, message.MessageID, string(messageType), textContent, caption,
		fileMetadata.FileID, fileMetadata.FileName, fileMetadata.FileSize, fileMetadata.MimeType, fileMetadata.Duration,
		forwardedDate, forwardedFrom,
		"{"+strings.Join(urls, ",")+"}",
		"{"+strings.Join(hashtags, ",")+"}",
		"{"+strings.Join(mentions, ",")+"}")

	return err
}

func generateForwardedTimes(message *tgbotapi.Message) (*time.Time, *string) {
	var forwardedDate *time.Time
	var forwardedFrom *string
	if message.ForwardFrom != nil {
		if message.ForwardDate != 0 {
			date := time.Unix(int64(message.ForwardDate), 0)
			forwardedDate = &date
		}
		from := message.ForwardFrom.FirstName
		if message.ForwardFrom.LastName != "" {
			from += " " + message.ForwardFrom.LastName
		}
		if message.ForwardFrom.UserName != "" {
			from += " (@" + message.ForwardFrom.UserName + ")"
		}
		forwardedFrom = &from
	}
	return forwardedDate, forwardedFrom
}
