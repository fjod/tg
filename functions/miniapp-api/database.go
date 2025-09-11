package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type User struct {
	ID         int64     `json:"id" db:"id"`
	TelegramID int64     `json:"telegram_id" db:"telegram_id"`
	Username   *string   `json:"username" db:"username"`
	FirstName  *string   `json:"first_name" db:"first_name"`
	LastName   *string   `json:"last_name" db:"last_name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	IsActive   bool      `json:"is_active" db:"is_active"`
}

type Tag struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	Name         string    `json:"name" db:"name"`
	Color        *string   `json:"color" db:"color"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	MessageCount int       `json:"message_count" db:"message_count"`
}

type MessageResponse struct {
	ID                int64     `json:"id" db:"id"`
	TelegramMessageID int64     `json:"telegram_message_id" db:"telegram_message_id"`
	MessageType       string    `json:"message_type" db:"message_type"`
	TextContent       *string   `json:"text_content" db:"text_content"`
	Caption           *string   `json:"caption" db:"caption"`
	FileName          *string   `json:"file_name" db:"file_name"`
	FileSize          *int64    `json:"file_size" db:"file_size"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	ForwardedFrom     *string   `json:"forwarded_from" db:"forwarded_from"`
	URLs              []string  `json:"urls"`
	Hashtags          []string  `json:"hashtags"`
}

func initDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL  environment variable not set")
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

func getUserTagsWithCounts(db *sql.DB, userID int64) ([]Tag, error) {
	query := `
		SELECT t.id, t.user_id, t.name, t.color, t.created_at, COUNT(mt.message_id) as message_count
		FROM tags t
		LEFT JOIN message_tags mt ON t.id = mt.tag_id
		WHERE t.user_id = $1
		GROUP BY t.id, t.user_id, t.name, t.color, t.created_at
		ORDER BY message_count DESC, t.name ASC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		var color sql.NullString

		if err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name, &color, &tag.CreatedAt, &tag.MessageCount); err != nil {
			return nil, err
		}

		if color.Valid {
			tag.Color = &color.String
		}

		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func getTagMessages(db *sql.DB, userID int64, tagID int64) ([]MessageResponse, error) {
	// First verify that the tag belongs to the user
	var tagExists bool
	tagQuery := "SELECT EXISTS(SELECT 1 FROM tags WHERE id = $1 AND user_id = $2)"
	err := db.QueryRow(tagQuery, tagID, userID).Scan(&tagExists)
	if err != nil {
		return nil, fmt.Errorf("failed to verify tag ownership: %v", err)
	}
	if !tagExists {
		return nil, fmt.Errorf("tag not found or access denied")
	}

	// Query messages for the specified tag
	query := `
		SELECT 
			m.id, 
			m.telegram_message_id, 
			m.message_type, 
			m.text_content, 
			m.caption, 
			m.file_name, 
			m.file_size, 
			m.created_at, 
			m.forwarded_from, 
			m.urls, 
			m.hashtags
		FROM messages m
		INNER JOIN message_tags mt ON m.id = mt.message_id
		WHERE mt.tag_id = $1 AND m.user_id = $2
		ORDER BY m.created_at DESC`

	rows, err := db.Query(query, tagID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %v", err)
	}
	defer rows.Close()

	var messages []MessageResponse
	for rows.Next() {
		var msg MessageResponse
		var textContent, caption, fileName, forwardedFrom sql.NullString
		var fileSize sql.NullInt64
		var urls, hashtags pq.StringArray

		err := rows.Scan(
			&msg.ID,
			&msg.TelegramMessageID,
			&msg.MessageType,
			&textContent,
			&caption,
			&fileName,
			&fileSize,
			&msg.CreatedAt,
			&forwardedFrom,
			&urls,
			&hashtags,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message row: %v", err)
		}

		// Handle nullable fields
		if textContent.Valid {
			msg.TextContent = &textContent.String
		}
		if caption.Valid {
			msg.Caption = &caption.String
		}
		if fileName.Valid {
			msg.FileName = &fileName.String
		}
		if fileSize.Valid {
			msg.FileSize = &fileSize.Int64
		}
		if forwardedFrom.Valid {
			msg.ForwardedFrom = &forwardedFrom.String
		}

		// Handle arrays (they might be nil, that's fine)
		msg.URLs = []string(urls)
		msg.Hashtags = []string(hashtags)

		// Ensure arrays are not nil for JSON serialization
		if msg.URLs == nil {
			msg.URLs = []string{}
		}
		if msg.Hashtags == nil {
			msg.Hashtags = []string{}
		}

		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
