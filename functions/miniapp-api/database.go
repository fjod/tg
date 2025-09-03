package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

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
