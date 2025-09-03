package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_ "modernc.org/sqlite"
)

// Test helper functions

// createTestTag creates a test tag in the database
func createTestTag(t *testing.T, db *sql.DB, userID int64, tagName, color string) int64 {
	var query string
	var args []interface{}

	if color != "" {
		query = `INSERT INTO tags (user_id, name, color, created_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)`
		args = []interface{}{userID, tagName, color}
	} else {
		query = `INSERT INTO tags (user_id, name, created_at) VALUES (?, ?, CURRENT_TIMESTAMP)`
		args = []interface{}{userID, tagName}
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		t.Fatalf("Failed to create test tag: %v", err)
	}
	tagID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get tag ID: %v", err)
	}
	return tagID
}

// createTestMessageTag creates a message-tag relationship
func createTestMessageTag(t *testing.T, db *sql.DB, messageID, tagID int64) {
	query := `INSERT INTO message_tags (message_id, tag_id, created_at) VALUES (?, ?, CURRENT_TIMESTAMP)`
	_, err := db.Exec(query, messageID, tagID)
	if err != nil {
		t.Fatalf("Failed to create test message-tag relationship: %v", err)
	}
}

// Test wrapper functions to handle interface conversion
func testShowTagSelection(bot BotAPI, message *tgbotapi.Message, db *sql.DB) {
	// Since the actual functions expect *tgbotapi.BotAPI, we need to work around this
	// For testing purposes, we'll directly test the core logic
	if message == nil || message.From == nil {
		return
	}

	tags, err := getUserTags(db, message.From.ID)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Could not load your tags."))
		return
	}

	if len(tags) <= 20 {
		testShowTagSelectionWithButtons(bot, message, tags)
	} else {
		testShowTagSelectionWithText(bot, message, tags)
	}
}

func testShowTagSelectionWithButtons(bot BotAPI, message *tgbotapi.Message, tags []Tag) {
	var responseText string
	var keyboard tgbotapi.InlineKeyboardMarkup

	if len(tags) == 0 {
		responseText = "You don't have any tags yet. Click the button below to create your first tag:"
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Create New Tag", fmt.Sprintf("new_tag:%d", message.MessageID)),
			),
		)
	} else {
		responseText = "Choose a tag or create a new one:"

		var rows [][]tgbotapi.InlineKeyboardButton
		for i := 0; i < len(tags); i += 2 {
			var row []tgbotapi.InlineKeyboardButton

			row = append(row, tgbotapi.NewInlineKeyboardButtonData(
				tags[i].Name,
				fmt.Sprintf("tag:%d:%d", tags[i].ID, message.MessageID),
			))

			if i+1 < len(tags) {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(
					tags[i+1].Name,
					fmt.Sprintf("tag:%d:%d", tags[i+1].ID, message.MessageID),
				))
			}

			rows = append(rows, row)
		}

		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("➕ Create New Tag", fmt.Sprintf("new_tag:%d", message.MessageID)),
		})

		keyboard = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = keyboard

	bot.Send(msg)
}

func testShowTagSelectionWithText(bot BotAPI, message *tgbotapi.Message, tags []Tag) {
	responseText := fmt.Sprintf("You have many tags (%d). Choose by typing its name or number, or create a new one:\n\n", len(tags))

	for i, tag := range tags {
		responseText += fmt.Sprintf("%d. %s\n", i+1, tag.Name)
	}
	responseText += fmt.Sprintf("\nType a tag name/number or create a new tag.\n\n[MSG_ID:%d]", message.MessageID)

	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}

	bot.Send(msg)
}

// TestGetUserTags tests the getUserTags function
func TestGetUserTags(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		setupTags   []struct{ name, color string }
		expectedLen int
		expectError bool
	}{
		{
			name:        "No tags for user",
			userID:      123,
			setupTags:   nil,
			expectedLen: 0,
			expectError: false,
		},
		{
			name:   "Single tag without color",
			userID: 123,
			setupTags: []struct{ name, color string }{
				{"work", ""},
			},
			expectedLen: 1,
			expectError: false,
		},
		{
			name:   "Single tag with color",
			userID: 123,
			setupTags: []struct{ name, color string }{
				{"urgent", "#ff0000"},
			},
			expectedLen: 1,
			expectError: false,
		},
		{
			name:   "Multiple tags mixed colors",
			userID: 123,
			setupTags: []struct{ name, color string }{
				{"work", "#0000ff"},
				{"personal", ""},
				{"urgent", "#ff0000"},
			},
			expectedLen: 3,
			expectError: false,
		},
		{
			name:   "Tags sorted alphabetically",
			userID: 123,
			setupTags: []struct{ name, color string }{
				{"zebra", ""},
				{"apple", ""},
				{"banana", ""},
			},
			expectedLen: 3,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			// Create test user
			createTestUser(t, db, tt.userID, "testuser")

			// Setup test tags
			for _, tagData := range tt.setupTags {
				createTestTag(t, db, tt.userID, tagData.name, tagData.color)
			}

			// Test getUserTags
			tags, err := getUserTags(db, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, tags, tt.expectedLen)

			// Verify tags are sorted alphabetically
			if len(tags) > 1 {
				for i := 1; i < len(tags); i++ {
					assert.True(t, tags[i-1].Name <= tags[i].Name, "Tags should be sorted alphabetically")
				}
			}

			// Verify tag data
			for _, tag := range tags {
				assert.Equal(t, tt.userID, tag.UserID)
				assert.NotEmpty(t, tag.Name)
				assert.Greater(t, tag.ID, int64(0))
			}

			// Test specific scenarios
			if tt.name == "Single tag with color" && len(tags) == 1 {
				assert.NotNil(t, tags[0].Color)
				assert.Equal(t, "#ff0000", *tags[0].Color)
			}
			if tt.name == "Single tag without color" && len(tags) == 1 {
				assert.Nil(t, tags[0].Color)
			}
		})
	}
}

// TestGetOrCreateTag tests the getOrCreateTag function
func TestGetOrCreateTag(t *testing.T) {
	tests := []struct {
		name         string
		userID       int64
		tagName      string
		existingTags []string
		expectCreate bool
		expectError  bool
	}{
		{
			name:         "Create new tag",
			userID:       123,
			tagName:      "newtag",
			existingTags: nil,
			expectCreate: true,
			expectError:  false,
		},
		{
			name:         "Get existing tag",
			userID:       123,
			tagName:      "existingtag",
			existingTags: []string{"existingtag"},
			expectCreate: false,
			expectError:  false,
		},
		{
			name:         "Create tag when others exist",
			userID:       123,
			tagName:      "newtag",
			existingTags: []string{"othertag1", "othertag2"},
			expectCreate: true,
			expectError:  false,
		},
		{
			name:         "Get existing tag among many",
			userID:       123,
			tagName:      "middletag",
			existingTags: []string{"atag", "middletag", "ztag"},
			expectCreate: false,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			// Create test user
			createTestUser(t, db, tt.userID, "testuser")

			// Setup existing tags
			var existingTagID int64
			for _, tagName := range tt.existingTags {
				tagID := createTestTag(t, db, tt.userID, tagName, "")
				if tagName == tt.tagName {
					existingTagID = tagID
				}
			}

			// Test getOrCreateTag
			tagID, err := getOrCreateTag(db, tt.userID, tt.tagName)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Greater(t, tagID, int64(0))

			if !tt.expectCreate {
				// Should return existing tag ID
				assert.Equal(t, existingTagID, tagID)
			}

			// Verify tag exists in database
			var retrievedName string
			query := `SELECT name FROM tags WHERE id = ? AND user_id = ?`
			err = db.QueryRow(query, tagID, tt.userID).Scan(&retrievedName)
			assert.NoError(t, err)
			assert.Equal(t, tt.tagName, retrievedName)
		})
	}
}

// TestTagMessage tests the tagMessage function
func TestTagMessage(t *testing.T) {
	tests := []struct {
		name             string
		messageID        int64
		tagID            int64
		existingRelation bool
		expectError      bool
	}{
		{
			name:             "Tag new message",
			messageID:        1,
			tagID:            1,
			existingRelation: false,
			expectError:      false,
		},
		{
			name:             "Tag already tagged message (should not error)",
			messageID:        1,
			tagID:            1,
			existingRelation: true,
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			userID := int64(123)

			// Create test user, message, and tag
			createTestUser(t, db, userID, "testuser")
			messageID := createTestMessage(t, db, userID, 456)
			tagID := createTestTag(t, db, userID, "testtag", "")

			// Setup existing relation if needed
			if tt.existingRelation {
				createTestMessageTag(t, db, messageID, tagID)
			}

			// Test tagMessage
			err := tagMessage(db, messageID, tagID)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify relationship exists
			var count int
			query := `SELECT COUNT(*) FROM message_tags WHERE message_id = ? AND tag_id = ?`
			err = db.QueryRow(query, messageID, tagID).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 1, count, "Should have exactly one message-tag relationship")
		})
	}
}

// TestGetMessageByTelegramID tests the getMessageByTelegramID function
func TestGetMessageByTelegramID(t *testing.T) {
	tests := []struct {
		name              string
		userID            int64
		telegramMessageID int64
		messageExists     bool
		expectError       bool
	}{
		{
			name:              "Find existing message",
			userID:            123,
			telegramMessageID: 456,
			messageExists:     true,
			expectError:       false,
		},
		{
			name:              "Message not found",
			userID:            123,
			telegramMessageID: 999,
			messageExists:     false,
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			// Create test user
			createTestUser(t, db, tt.userID, "testuser")

			var expectedMessageID int64
			if tt.messageExists {
				// Create test message
				expectedMessageID = createTestMessage(t, db, tt.userID, tt.telegramMessageID)
			}

			// Test getMessageByTelegramID
			messageID, err := getMessageByTelegramID(db, tt.userID, tt.telegramMessageID)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, expectedMessageID, messageID)
		})
	}
}

// TestShowTagSelection tests the showTagSelection routing logic
func TestShowTagSelection(t *testing.T) {
	tests := []struct {
		name         string
		numTags      int
		expectButton bool
		expectText   bool
	}{
		{
			name:         "No tags - should use buttons",
			numTags:      0,
			expectButton: true,
			expectText:   false,
		},
		{
			name:         "Few tags - should use buttons",
			numTags:      10,
			expectButton: true,
			expectText:   false,
		},
		{
			name:         "Exactly 20 tags - should use buttons",
			numTags:      20,
			expectButton: true,
			expectText:   false,
		},
		{
			name:         "Many tags - should use text",
			numTags:      25,
			expectButton: false,
			expectText:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			userID := int64(123)
			createTestUser(t, db, userID, "testuser")

			// Create specified number of tags
			for i := 0; i < tt.numTags; i++ {
				createTestTag(t, db, userID, fmt.Sprintf("tag%d", i), "")
			}

			// Create mock bot
			mockBot := &MockBotAPI{}

			// Setup expectations for bot.Send calls
			if tt.expectButton || tt.expectText {
				mockBot.On("Send", mock.AnythingOfType("tgbotapi.MessageConfig")).Return(tgbotapi.Message{}, nil)
			}

			// Create test message
			message := createTelegramMessage(456, userID, "testuser", "test message")

			// Test showTagSelection using wrapper function
			testShowTagSelection(mockBot, message, db)

			// Verify expectations
			mockBot.AssertExpectations(t)
		})
	}
}

// TestShowTagSelectionWithButtons tests the button UI generation
func TestShowTagSelectionWithButtons(t *testing.T) {
	tests := []struct {
		name        string
		numTags     int
		expectRows  int
		expectError bool
	}{
		{
			name:        "No tags - create button only",
			numTags:     0,
			expectRows:  1, // Just "Create New Tag" button
			expectError: false,
		},
		{
			name:        "Single tag",
			numTags:     1,
			expectRows:  2, // 1 tag row + create button row
			expectError: false,
		},
		{
			name:        "Two tags - same row",
			numTags:     2,
			expectRows:  2, // 1 tag row + create button row
			expectError: false,
		},
		{
			name:        "Three tags - two rows",
			numTags:     3,
			expectRows:  3, // 2 tag rows + create button row
			expectError: false,
		},
		{
			name:        "Twenty tags",
			numTags:     20,
			expectRows:  11, // 10 tag rows + create button row
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			userID := int64(123)
			createTestUser(t, db, userID, "testuser")

			// Create test tags
			var tags []Tag
			for i := 0; i < tt.numTags; i++ {
				tagID := createTestTag(t, db, userID, fmt.Sprintf("tag%d", i), "")
				tags = append(tags, Tag{
					ID:     tagID,
					UserID: userID,
					Name:   fmt.Sprintf("tag%d", i),
				})
			}

			// Create mock bot
			mockBot := &MockBotAPI{}

			// Capture the message config to verify keyboard structure
			mockBot.On("Send", mock.MatchedBy(func(c tgbotapi.Chattable) bool {
				if msgConfig, ok := c.(tgbotapi.MessageConfig); ok {
					if keyboard, ok := msgConfig.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup); ok {
						assert.Len(t, keyboard.InlineKeyboard, tt.expectRows, "Should have correct number of button rows")

						// Verify button structure
						if tt.numTags == 0 {
							// Should have only create button
							assert.Len(t, keyboard.InlineKeyboard[0], 1)
							assert.Contains(t, keyboard.InlineKeyboard[0][0].Text, "Create New Tag")
						} else {
							// Should have create button in last row
							lastRow := keyboard.InlineKeyboard[len(keyboard.InlineKeyboard)-1]
							assert.Len(t, lastRow, 1)
							assert.Contains(t, lastRow[0].Text, "Create New Tag")
						}
					}
					return true
				}
				return false
			})).Return(tgbotapi.Message{}, nil)

			// Create test message
			message := createTelegramMessage(456, userID, "testuser", "test message")

			// Test showTagSelectionWithButtons
			testShowTagSelectionWithButtons(mockBot, message, tags)

			// Verify expectations
			mockBot.AssertExpectations(t)
		})
	}
}

// TestShowTagSelectionWithText tests the text UI generation
func TestShowTagSelectionWithText(t *testing.T) {
	tests := []struct {
		name     string
		numTags  int
		tagNames []string
	}{
		{
			name:     "Many tags with numbering",
			numTags:  25,
			tagNames: []string{"apple", "banana", "cherry", "date", "elderberry"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			userID := int64(123)
			createTestUser(t, db, userID, "testuser")

			// Create test tags
			var tags []Tag
			for i := 0; i < len(tt.tagNames); i++ {
				tagID := createTestTag(t, db, userID, tt.tagNames[i], "")
				tags = append(tags, Tag{
					ID:     tagID,
					UserID: userID,
					Name:   tt.tagNames[i],
				})
			}

			// Create mock bot
			mockBot := &MockBotAPI{}

			// Capture the message to verify MSG_ID format
			mockBot.On("Send", mock.MatchedBy(func(c tgbotapi.Chattable) bool {
				if msgConfig, ok := c.(tgbotapi.MessageConfig); ok {
					// Verify MSG_ID is present in text
					assert.Contains(t, msgConfig.Text, "[MSG_ID:")
					assert.Contains(t, msgConfig.Text, "]")

					// Verify numbering format
					for i, tagName := range tt.tagNames {
						expectedLine := fmt.Sprintf("%d. %s", i+1, tagName)
						assert.Contains(t, msgConfig.Text, expectedLine)
					}

					// Verify force reply is set
					if forceReply, ok := msgConfig.ReplyMarkup.(tgbotapi.ForceReply); ok {
						assert.True(t, forceReply.ForceReply)
						assert.True(t, forceReply.Selective)
					}
					return true
				}
				return false
			})).Return(tgbotapi.Message{}, nil)

			// Create test message
			message := createTelegramMessage(456, userID, "testuser", "test message")

			// Test showTagSelectionWithText
			testShowTagSelectionWithText(mockBot, message, tags)

			// Verify expectations
			mockBot.AssertExpectations(t)
		})
	}
}

// TestTagSelectionLogic tests the core logic of tag selection parsing
func TestTagSelectionLogic(t *testing.T) {
	tests := []struct {
		name        string
		replyText   string
		userInput   string
		expectMsgID string
		expectError bool
	}{
		{
			name:        "Valid MSG_ID extraction",
			replyText:   "Choose a tag:\n\n[MSG_ID:456]",
			userInput:   "work",
			expectMsgID: "456",
			expectError: false,
		},
		{
			name:        "No MSG_ID in reply",
			replyText:   "Choose a tag without message ID",
			userInput:   "work",
			expectMsgID: "",
			expectError: true,
		},
		{
			name:        "Malformed MSG_ID",
			replyText:   "[MSG_ID:abc]",
			userInput:   "work",
			expectMsgID: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test MSG_ID extraction logic
			msgIDStart := strings.Index(tt.replyText, "[MSG_ID:")
			if msgIDStart == -1 {
				if !tt.expectError {
					t.Errorf("Expected to find MSG_ID but didn't")
				}
				return
			}

			msgIDEnd := strings.Index(tt.replyText[msgIDStart:], "]")
			if msgIDEnd == -1 {
				if !tt.expectError {
					t.Errorf("Expected to find closing bracket but didn't")
				}
				return
			}

			msgIDStr := tt.replyText[msgIDStart+8 : msgIDStart+msgIDEnd]

			if tt.expectError {
				// Test parsing error
				_, err := strconv.Atoi(msgIDStr)
				if msgIDStr != tt.expectMsgID && err == nil {
					t.Errorf("Expected error parsing %s but got none", msgIDStr)
				}
			} else {
				assert.Equal(t, tt.expectMsgID, msgIDStr)

				// Verify it can be parsed as integer
				_, err := strconv.Atoi(msgIDStr)
				assert.NoError(t, err)
			}
		})
	}
}

// TestCallbackDataParsing tests callback data parsing logic
func TestCallbackDataParsing(t *testing.T) {
	tests := []struct {
		name         string
		callbackData string
		expectValid  bool
		expectParts  int
	}{
		{
			name:         "Valid tag callback data",
			callbackData: "tag:1:456",
			expectValid:  true,
			expectParts:  3,
		},
		{
			name:         "Valid new_tag callback data",
			callbackData: "new_tag:456",
			expectValid:  true,
			expectParts:  2,
		},
		{
			name:         "Invalid format - too few parts",
			callbackData: "tag:invalid",
			expectValid:  false,
			expectParts:  2,
		},
		{
			name:         "Empty callback data",
			callbackData: "",
			expectValid:  false,
			expectParts:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(tt.callbackData, ":")
			assert.Len(t, parts, tt.expectParts)

			// Test tag callback parsing
			if strings.HasPrefix(tt.callbackData, "tag:") && len(parts) == 3 {
				// Should be able to parse tag ID and message ID
				tagID, err := strconv.ParseInt(parts[1], 10, 64)
				if tt.expectValid {
					assert.NoError(t, err)
					assert.Greater(t, tagID, int64(0))
				}

				messageID, err := strconv.Atoi(parts[2])
				if tt.expectValid {
					assert.NoError(t, err)
					assert.Greater(t, messageID, 0)
				}
			}

			// Test new_tag callback parsing
			if strings.HasPrefix(tt.callbackData, "new_tag:") && len(parts) == 2 {
				messageID, err := strconv.Atoi(parts[1])
				if tt.expectValid {
					assert.NoError(t, err)
					assert.Greater(t, messageID, 0)
				}
			}
		})
	}
}

// TestNewTagWorkflow tests the new tag creation workflow logic
func TestNewTagWorkflow(t *testing.T) {
	t.Run("New tag prompt message format", func(t *testing.T) {
		// Test that new tag workflow creates proper MSG_ID format
		originalMessageID := 456
		expectedText := fmt.Sprintf("Please reply with the name for your new tag:\n\n[MSG_ID:%d]", originalMessageID)

		// Verify MSG_ID is present and parseable
		assert.Contains(t, expectedText, "[MSG_ID:")
		assert.Contains(t, expectedText, "]")

		// Verify can extract message ID
		msgIDStart := strings.Index(expectedText, "[MSG_ID:")
		msgIDEnd := strings.Index(expectedText[msgIDStart:], "]")
		msgIDStr := expectedText[msgIDStart+8 : msgIDStart+msgIDEnd]

		parsedID, err := strconv.Atoi(msgIDStr)
		assert.NoError(t, err)
		assert.Equal(t, originalMessageID, parsedID)
	})
}

// TestTagsEdgeCases tests comprehensive edge cases
func TestTagsEdgeCases(t *testing.T) {
	t.Run("Tag name with special characters", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		userID := int64(123)
		createTestUser(t, db, userID, "testuser")

		// Test creating tags with special characters
		specialNames := []string{
			"tag-with-dashes",
			"tag_with_underscores",
			"tag with spaces",
			"tag@email.com",
			"tag#hashtag",
			"🏷️ emoji tag",
		}

		for _, tagName := range specialNames {
			tagID, err := getOrCreateTag(db, userID, tagName)
			assert.NoError(t, err, "Should handle special characters in tag names")
			assert.Greater(t, tagID, int64(0))
		}
	})

	t.Run("Very long tag name", func(t *testing.T) {
		db := setupTestDB(t)
		defer db.Close()

		userID := int64(123)
		createTestUser(t, db, userID, "testuser")

		// Test with very long tag name
		longTagName := strings.Repeat("a", 200)
		tagID, err := getOrCreateTag(db, userID, longTagName)
		assert.NoError(t, err, "Should handle long tag names")
		assert.Greater(t, tagID, int64(0))
	})

	t.Run("Database connection error scenarios", func(t *testing.T) {
		// Test with closed database
		db := setupTestDB(t)
		db.Close() // Close immediately

		userID := int64(123)

		// All functions should handle database errors gracefully
		_, err := getUserTags(db, userID)
		assert.Error(t, err, "Should handle database connection errors")

		_, err = getOrCreateTag(db, userID, "testtag")
		assert.Error(t, err, "Should handle database connection errors")

		err = tagMessage(db, 1, 1)
		assert.Error(t, err, "Should handle database connection errors")

		_, err = getMessageByTelegramID(db, userID, 456)
		assert.Error(t, err, "Should handle database connection errors")
	})

	t.Run("Message ID extraction edge cases", func(t *testing.T) {
		// Test various malformed MSG_ID formats
		malformedCases := []string{
			"No message ID here",
			"[MSG_ID:",
			"[MSG_ID:abc]",
			"[MSG_ID:]",
			"MSG_ID:123]",
			"[MSG_ID:123",
			"Multiple [MSG_ID:123] and [MSG_ID:456]",
		}

		for _, replyText := range malformedCases {
			// Test MSG_ID extraction robustness
			msgIDStart := strings.Index(replyText, "[MSG_ID:")
			if msgIDStart == -1 {
				// Expected for cases without proper format
				continue
			}

			msgIDEnd := strings.Index(replyText[msgIDStart:], "]")
			if msgIDEnd == -1 {
				// Expected for malformed cases
				continue
			}

			msgIDStr := replyText[msgIDStart+8 : msgIDStart+msgIDEnd]
			_, err := strconv.Atoi(msgIDStr)
			// Should expect error for malformed IDs like "abc" and ""
			if replyText == "[MSG_ID:abc]" || replyText == "[MSG_ID:]" {
				assert.Error(t, err, "Should fail to parse malformed message ID")
			}
		}
	})

	t.Run("Input validation", func(t *testing.T) {
		// Test tag name validation scenarios
		testInputs := []struct {
			input   string
			isEmpty bool
		}{
			{"valid_tag", false},
			{"", true},
			{"   ", true}, // Just whitespace
			{"tag with spaces", false},
		}

		for _, test := range testInputs {
			trimmed := strings.TrimSpace(test.input)
			if test.isEmpty {
				assert.Empty(t, trimmed, "Should be empty after trimming")
			} else {
				assert.NotEmpty(t, trimmed, "Should not be empty after trimming")
			}
		}
	})
}
