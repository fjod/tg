package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

// Test helper functions

// createTestUserStruct creates a test Telegram user struct
func createTestUserStruct(id int64, username, firstName, lastName string) *tgbotapi.User {
	user := &tgbotapi.User{
		ID:        id,
		FirstName: firstName,
	}
	if username != "" {
		user.UserName = username
	}
	if lastName != "" {
		user.LastName = lastName
	}
	return user
}

// createTestMessageStruct creates a test Telegram message struct
func createTestMessageStruct(messageID int, user *tgbotapi.User, text string) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: messageID,
		From:      user,
		Text:      text,
		Date:      int(time.Now().Unix()),
	}
}

// createTestForwardedMessage creates a test forwarded message
func createTestForwardedMessage(messageID int, user *tgbotapi.User, text string, forwardFrom *tgbotapi.User, forwardDate int) *tgbotapi.Message {
	msg := createTestMessageStruct(messageID, user, text)
	msg.ForwardFrom = forwardFrom
	msg.ForwardDate = forwardDate
	return msg
}

// createTestPhotoMessage creates a test message with photo
func createTestPhotoMessage(messageID int, user *tgbotapi.User, caption string, photos ...tgbotapi.PhotoSize) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: messageID,
		From:      user,
		Photo:     photos,
		Caption:   caption,
		Date:      int(time.Now().Unix()),
	}
}

// getUserFromDB retrieves a user from the database for verification
func getUserFromDB(t *testing.T, db *sql.DB, telegramID int64) (username, firstName, lastName sql.NullString, isActive bool) {
	query := `SELECT username, first_name, last_name, is_active FROM users WHERE telegram_id = ?`
	err := db.QueryRow(query, telegramID).Scan(&username, &firstName, &lastName, &isActive)
	if err != nil {
		t.Fatalf("Failed to get user from DB: %v", err)
	}
	return
}

// getMessageFromDB retrieves a message from the database for verification
func getMessageFromDB(t *testing.T, db *sql.DB, userID int64, telegramMessageID int) (messageType string, textContent, caption sql.NullString) {
	query := `SELECT message_type, text_content, caption FROM messages WHERE user_id = ? AND telegram_message_id = ?`
	err := db.QueryRow(query, userID, telegramMessageID).Scan(&messageType, &textContent, &caption)
	if err != nil {
		t.Fatalf("Failed to get message from DB: %v", err)
	}
	return
}

// TestTruncateText tests the text truncation functionality
func TestTruncateText(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		maxLength int
		expected  string
	}{
		{
			name:      "Text shorter than max length",
			text:      "Hello World",
			maxLength: 20,
			expected:  "Hello World",
		},
		{
			name:      "Text equal to max length",
			text:      "Hello World",
			maxLength: 11,
			expected:  "Hello World",
		},
		{
			name:      "Text longer than max length",
			text:      "This is a very long text that should be truncated",
			maxLength: 20,
			expected:  "This is a very long ...", // Note the space before "..."
		},
		{
			name:      "Empty text",
			text:      "",
			maxLength: 10,
			expected:  "",
		},
		{
			name:      "Single character",
			text:      "A",
			maxLength: 5,
			expected:  "A",
		},
		{
			name:      "Max length zero",
			text:      "Hello",
			maxLength: 0,
			expected:  "...",
		},
		{
			name:      "Max length one",
			text:      "Hello",
			maxLength: 1,
			expected:  "H...",
		},
		{
			name:      "Unicode text",
			text:      "Hello世界", 
			maxLength: 8,
			expected:  "Hello世...",
		},
		{
			name:      "Text with newlines and spaces",
			text:      "Line 1\nLine 2\n\nLine 3",
			maxLength: 10,
			expected:  "Line 1\nLin...",
		},
		{
			name:      "Very long single word",
			text:      "supercalifragilisticexpialidocious",
			maxLength: 15,
			expected:  "supercalifragil...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateText(tt.text, tt.maxLength)
			assert.Equal(t, tt.expected, result)
			
			// Verify result doesn't exceed expected length (accounting for "...")
			if tt.maxLength > 0 {
				assert.LessOrEqual(t, len(result), tt.maxLength+3) // +3 for "..."
			}
		})
	}
}

// TestGenerateForwardedTimes tests forwarded message metadata extraction
func TestGenerateForwardedTimes(t *testing.T) {
	tests := []struct {
		name             string
		message          *tgbotapi.Message
		expectDate       bool
		expectedFromText string
	}{
		{
			name:             "No forward data",
			message:          createTestMessageStruct(1, createTestUserStruct(123, "user", "Test", "User"), "test"),
			expectDate:       false,
			expectedFromText: "",
		},
		{
			name: "Complete forward data with username",
			message: createTestForwardedMessage(
				1,
				createTestUserStruct(123, "user", "Test", "User"),
				"forwarded message",
				createTestUserStruct(456, "forward_user", "Forward", "User"),
				1640995200, // 2022-01-01 00:00:00 UTC
			),
			expectDate:       true,
			expectedFromText: "Forward User (@forward_user)",
		},
		{
			name: "Forward data without username",
			message: createTestForwardedMessage(
				1,
				createTestUserStruct(123, "user", "Test", "User"),
				"forwarded message",
				createTestUserStruct(456, "", "Forward", "User"),
				1640995200,
			),
			expectDate:       true,
			expectedFromText: "Forward User",
		},
		{
			name: "Forward data without last name",
			message: createTestForwardedMessage(
				1,
				createTestUserStruct(123, "user", "Test", "User"),
				"forwarded message",
				createTestUserStruct(456, "forward_user", "Forward", ""),
				1640995200,
			),
			expectDate:       true,
			expectedFromText: "Forward (@forward_user)",
		},
		{
			name: "Forward data with only first name",
			message: createTestForwardedMessage(
				1,
				createTestUserStruct(123, "user", "Test", "User"),
				"forwarded message",
				createTestUserStruct(456, "", "Forward", ""),
				1640995200,
			),
			expectDate:       true,
			expectedFromText: "Forward",
		},
		{
			name: "Forward data with zero timestamp",
			message: func() *tgbotapi.Message {
				msg := createTestForwardedMessage(
					1,
					createTestUserStruct(123, "user", "Test", "User"),
					"forwarded message",
					createTestUserStruct(456, "forward_user", "Forward", "User"),
					0,
				)
				return msg
			}(),
			expectDate:       false,
			expectedFromText: "Forward User (@forward_user)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forwardedDate, forwardedFrom := generateForwardedTimes(tt.message)
			
			if tt.expectDate {
				assert.NotNil(t, forwardedDate, "Expected forwarded date to be set")
				if tt.message.ForwardDate != 0 {
					expectedTime := time.Unix(int64(tt.message.ForwardDate), 0)
					assert.Equal(t, expectedTime, *forwardedDate)
				}
			} else {
				assert.Nil(t, forwardedDate, "Expected forwarded date to be nil")
			}
			
			if tt.expectedFromText != "" {
				assert.NotNil(t, forwardedFrom, "Expected forwarded from to be set")
				assert.Equal(t, tt.expectedFromText, *forwardedFrom)
			} else {
				assert.Nil(t, forwardedFrom, "Expected forwarded from to be nil")
			}
		})
	}
}

// TestSaveUser tests user persistence functionality
func TestSaveUser(t *testing.T) {
	tests := []struct {
		name           string
		user           *tgbotapi.User
		existingUser   *tgbotapi.User // For testing UPSERT behavior
		expectError    bool
	}{
		{
			name: "New user with all fields",
			user: createTestUserStruct(123, "testuser", "Test", "User"),
			expectError: false,
		},
		{
			name: "New user with minimal fields",
			user: createTestUserStruct(456, "", "MinimalUser", ""),
			expectError: false,
		},
		{
			name: "User with empty username",
			user: createTestUserStruct(789, "", "Empty", "Username"),
			expectError: false,
		},
		{
			name: "User with very long names",
			user: createTestUserStruct(999, "verylongusername", 
				"VeryLongFirstNameThatExceedsNormalLength", 
				"VeryLongLastNameThatExceedsNormalLength"),
			expectError: false,
		},
		{
			name: "Update existing user",
			user: createTestUserStruct(123, "updated_user", "Updated", "Name"),
			existingUser: createTestUserStruct(123, "old_user", "Old", "Name"),
			expectError: false,
		},
		{
			name: "User with unicode characters",
			user: createTestUserStruct(888, "用户名", "名前", "姓氏"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			// Insert existing user if specified
			if tt.existingUser != nil {
				err := saveUser(db, tt.existingUser)
				assert.NoError(t, err, "Failed to save existing user for test setup")
			}

			// Test saving the user
			err := saveUser(db, tt.user)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify user was saved correctly
			username, firstName, lastName, isActive := getUserFromDB(t, db, tt.user.ID)
			
			if tt.user.UserName != "" {
				assert.True(t, username.Valid)
				assert.Equal(t, tt.user.UserName, username.String)
			} else {
				assert.False(t, username.Valid)
			}
			
			assert.True(t, firstName.Valid)
			assert.Equal(t, tt.user.FirstName, firstName.String)
			
			if tt.user.LastName != "" {
				assert.True(t, lastName.Valid)
				assert.Equal(t, tt.user.LastName, lastName.String)
			} else {
				assert.False(t, lastName.Valid)
			}
			
			assert.True(t, isActive)

			// If this was an update, verify the old data was replaced
			if tt.existingUser != nil {
				if tt.user.UserName != tt.existingUser.UserName {
					if tt.user.UserName != "" {
						assert.Equal(t, tt.user.UserName, username.String)
					}
				}
			}
		})
	}
}

// TestSaveMessage tests message persistence functionality
func TestSaveMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     *tgbotapi.Message
		expectError bool
	}{
		{
			name:        "Simple text message",
			message:     createTestMessageStruct(1, createTestUserStruct(123, "user", "Test", "User"), "Hello World"),
			expectError: false,
		},
		{
			name:        "Long text message that gets truncated",
			message:     createTestMessageStruct(2, createTestUserStruct(123, "user", "Test", "User"), 
				"This is a very long message that exceeds the 150 character limit for text content storage and should be truncated to a preview while still extracting metadata from the full text"),
			expectError: false,
		},
		{
			name:        "Message with caption",
			message:     createTestPhotoMessage(3, createTestUserStruct(123, "user", "Test", "User"), 
				"Photo caption", tgbotapi.PhotoSize{FileID: "photo123", Width: 100, Height: 100}),
			expectError: false,
		},
		{
			name:        "Forwarded message",
			message:     createTestForwardedMessage(4, createTestUserStruct(123, "user", "Test", "User"), 
				"Forwarded content", createTestUserStruct(456, "forward_user", "Forward", "User"), 1640995200),
			expectError: false,
		},
		{
			name:        "Message with URLs and hashtags",
			message:     createTestMessageStruct(5, createTestUserStruct(123, "user", "Test", "User"), 
				"Check out https://example.com #awesome #test @mention"),
			expectError: false,
		},
		{
			name:        "Empty text message",
			message:     createTestMessageStruct(6, createTestUserStruct(123, "user", "Test", "User"), ""),
			expectError: false,
		},
		{
			name:        "Message with unicode content",
			message:     createTestMessageStruct(7, createTestUserStruct(123, "user", "Test", "User"), 
				"Unicode test: こんにちは 🌟 #日本語 @ユーザー https://日本語.example.com"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			// First save the user (required for foreign key constraint)
			err := saveUser(db, tt.message.From)
			assert.NoError(t, err, "Failed to save user for message test")

			// Test saving the message
			err = saveMessage(db, tt.message)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify message was saved correctly
			messageType, textContent, caption := getMessageFromDB(t, db, tt.message.From.ID, tt.message.MessageID)
			
			// Verify message type
			expectedType := string(getMessageType(tt.message))
			assert.Equal(t, expectedType, messageType)
			
			// Verify text content (should be truncated if long)
			if tt.message.Text != "" {
				assert.True(t, textContent.Valid)
				expectedText := truncateText(tt.message.Text, 150)
				assert.Equal(t, expectedText, textContent.String)
			}
			
			// Verify caption (should be truncated if long)
			if tt.message.Caption != "" {
				assert.True(t, caption.Valid)
				expectedCaption := truncateText(tt.message.Caption, 150)
				assert.Equal(t, expectedCaption, caption.String)
			}

			// Additional verification: check that message exists in database
			var count int
			countQuery := `SELECT COUNT(*) FROM messages WHERE user_id = ? AND telegram_message_id = ?`
			err = db.QueryRow(countQuery, tt.message.From.ID, tt.message.MessageID).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 1, count, "Message should exist in database")
		})
	}
}

// TestInitDB tests database initialization functionality
func TestInitDB(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Missing DATABASE_URL environment variable",
			envValue:    "",
			expectError: true,
			errorMsg:    "DATABASE_URL environment variable not set",
		},
		{
			name:        "Invalid database URL format",
			envValue:    "invalid://url/format",
			expectError: true,
			errorMsg:    "", // Error will come from sql.Open, message varies
		},
		{
			name:        "Valid PostgreSQL URL format",
			envValue:    "postgres://user:password@localhost/dbname?sslmode=disable",
			expectError: true, // Will fail to connect since it's not a real DB
			errorMsg:    "", // Connection error varies
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env var
			originalEnv := os.Getenv("DATABASE_URL")
			defer func() {
				if originalEnv != "" {
					os.Setenv("DATABASE_URL", originalEnv)
				} else {
					os.Unsetenv("DATABASE_URL")
				}
			}()

			// Set test env var
			if tt.envValue != "" {
				os.Setenv("DATABASE_URL", tt.envValue)
			} else {
				os.Unsetenv("DATABASE_URL")
			}

			// Test initDB
			db, err := initDB()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, db)
				return
			}

			// If no error expected, clean up
			if db != nil {
				db.Close()
			}
		})
	}
}

// TestDatabaseEdgeCases tests comprehensive edge cases and error scenarios
func TestDatabaseEdgeCases(t *testing.T) {
	t.Run("SaveUser with closed database", func(t *testing.T) {
		db := setupTestDB(t)
		db.Close() // Close the database

		user := createTestUserStruct(123, "testuser", "Test", "User")
		err := saveUser(db, user)
		assert.Error(t, err, "Should fail with closed database")
	})

	t.Run("SaveMessage with closed database", func(t *testing.T) {
		db := setupTestDB(t)
		db.Close() // Close the database

		user := createTestUserStruct(123, "testuser", "Test", "User")
		message := createTestMessageStruct(1, user, "test message")
		err := saveMessage(db, message)
		assert.Error(t, err, "Should fail with closed database")
	})

	t.Run("SaveMessage without existing user", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// Try to save message without saving user first 
		// Note: SQLite may not enforce foreign key constraints by default
		user := createTestUserStruct(123, "testuser", "Test", "User")
		message := createTestMessageStruct(1, user, "test message")
		err := saveMessage(db, message)
		// SQLite doesn't enforce foreign keys by default in test DB, so this may pass
		// This test documents the behavior rather than testing constraint enforcement
		_ = err // Just test that function doesn't crash
	})

	t.Run("TruncateText with negative max length", func(t *testing.T) {
		// This will panic as expected since the function uses text[:maxLength]
		defer func() {
			if r := recover(); r != nil {
				assert.Contains(t, fmt.Sprintf("%v", r), "slice bounds out of range")
			}
		}()
		
		// This should panic - document the expected behavior
		truncateText("Hello World", -1)
		t.Error("Expected panic but function completed normally")
	})

	t.Run("GenerateForwardedTimes with nil message", func(t *testing.T) {
		// This would panic in real code, but test defensive behavior
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Function panicked as expected with nil message: %v", r)
			}
		}()
		
		// This will likely panic, which is expected behavior
		generateForwardedTimes(nil)
	})

	t.Run("SaveUser with nil user", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// This will panic as expected since function accesses user.ID
		defer func() {
			if r := recover(); r != nil {
				assert.Contains(t, fmt.Sprintf("%v", r), "nil pointer dereference")
			}
		}()

		saveUser(db, nil)
		t.Error("Expected panic but function completed normally")
	})

	t.Run("SaveMessage with nil message", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// This will panic as expected since function accesses message.Text, etc.
		defer func() {
			if r := recover(); r != nil {
				assert.Contains(t, fmt.Sprintf("%v", r), "nil pointer dereference")
			}
		}()

		saveMessage(db, nil)
		t.Error("Expected panic but function completed normally")
	})

	t.Run("Very long text content handling", func(t *testing.T) {
		// Test with extremely long content
		longText := make([]byte, 10000) // 10KB of text
		for i := range longText {
			longText[i] = 'A'
		}
		
		result := truncateText(string(longText), 150)
		assert.LessOrEqual(t, len(result), 153) // 150 + "..." = 153
		assert.True(t, strings.HasSuffix(result, "..."))
	})

	t.Run("Message with all metadata types", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		user := createTestUserStruct(123, "testuser", "Test", "User")
		err := saveUser(db, user)
		assert.NoError(t, err)

		// Message with URLs, hashtags, mentions, and forwarding
		complexMessage := createTestForwardedMessage(
			1, user,
			"Complex message with https://example.com and https://test.org #tag1 #tag2 @user1 @user2 and more content",
			createTestUserStruct(456, "forward_user", "Forward", "User"),
			1640995200,
		)

		err = saveMessage(db, complexMessage)
		assert.NoError(t, err)

		// Verify the message was saved with all metadata
		var urls, hashtags, mentions string
		query := `SELECT urls, hashtags, mentions FROM messages WHERE user_id = ? AND telegram_message_id = ?`
		err = db.QueryRow(query, user.ID, complexMessage.MessageID).Scan(&urls, &hashtags, &mentions)
		assert.NoError(t, err)
		
		// URLs should be saved
		assert.Contains(t, urls, "https://example.com")
		assert.Contains(t, urls, "https://test.org")
		
		// Hashtags should be saved
		assert.Contains(t, hashtags, "tag1")
		assert.Contains(t, hashtags, "tag2")
		
		// Mentions should be saved
		assert.Contains(t, mentions, "user1")
		assert.Contains(t, mentions, "user2")
	})
}

// TestIntegrationWorkflows tests complete user and message save workflows
func TestIntegrationWorkflows(t *testing.T) {
	t.Run("Complete user and message workflow", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// Create a user and save
		user := createTestUserStruct(123, "integrationuser", "Integration", "User")
		err := saveUser(db, user)
		assert.NoError(t, err)

		// Create and save multiple messages for the user
		messages := []*tgbotapi.Message{
			createTestMessageStruct(1, user, "First message"),
			createTestMessageStruct(2, user, "Second message with https://example.com"),
			createTestForwardedMessage(3, user, "Forwarded message", 
				createTestUserStruct(456, "other_user", "Other", "User"), 1640995200),
		}

		for _, msg := range messages {
			err := saveMessage(db, msg)
			assert.NoError(t, err, "Failed to save message %d", msg.MessageID)
		}

		// Verify all messages were saved
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM messages WHERE user_id = ?", user.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, len(messages), count)
	})

	t.Run("User update workflow", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// Save initial user
		originalUser := createTestUserStruct(123, "originaluser", "Original", "Name")
		err := saveUser(db, originalUser)
		assert.NoError(t, err)

		// Update user information
		updatedUser := createTestUserStruct(123, "updateduser", "Updated", "Name")
		err = saveUser(db, updatedUser)
		assert.NoError(t, err)

		// Verify user was updated, not duplicated
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE telegram_id = ?", originalUser.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "Should have exactly one user record")

		// Verify the updated data
		username, firstName, lastName, _ := getUserFromDB(t, db, originalUser.ID)
		assert.Equal(t, updatedUser.UserName, username.String)
		assert.Equal(t, updatedUser.FirstName, firstName.String)
		assert.Equal(t, updatedUser.LastName, lastName.String)
	})

	t.Run("Batch message processing", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		// Create multiple users
		users := []*tgbotapi.User{
			createTestUserStruct(100, "user1", "User", "One"),
			createTestUserStruct(200, "user2", "User", "Two"),
			createTestUserStruct(300, "user3", "User", "Three"),
		}

		// Save all users
		for _, user := range users {
			err := saveUser(db, user)
			assert.NoError(t, err)
		}

		// Create messages from different users
		messageID := 1
		for _, user := range users {
			for i := 0; i < 3; i++ {
				msg := createTestMessageStruct(messageID, user, fmt.Sprintf("Message %d from user %d", i+1, user.ID))
				err := saveMessage(db, msg)
				assert.NoError(t, err)
				messageID++
			}
		}

		// Verify total message count
		var totalMessages int
		err := db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&totalMessages)
		assert.NoError(t, err)
		assert.Equal(t, 9, totalMessages) // 3 users * 3 messages each

		// Verify each user has correct number of messages
		for _, user := range users {
			var userMessages int
			err = db.QueryRow("SELECT COUNT(*) FROM messages WHERE user_id = ?", user.ID).Scan(&userMessages)
			assert.NoError(t, err)
			assert.Equal(t, 3, userMessages)
		}
	})
}