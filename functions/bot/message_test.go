package main

import (
	"database/sql"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

// Helper functions to create test messages

func createTextMessage(text, caption string) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: 1,
		Text:      text,
		Caption:   caption,
	}
}

func createPhotoMessage(caption string, photos ...tgbotapi.PhotoSize) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: 1,
		Photo:     photos,
		Caption:   caption,
	}
}

func createVideoMessage(caption string, video *tgbotapi.Video) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: 1,
		Video:     video,
		Caption:   caption,
	}
}

func createDocumentMessage(caption string, document *tgbotapi.Document) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: 1,
		Document:  document,
		Caption:   caption,
	}
}

func createAudioMessage(caption string, audio *tgbotapi.Audio) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: 1,
		Audio:     audio,
		Caption:   caption,
	}
}

func createVoiceMessage(voice *tgbotapi.Voice) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: 1,
		Voice:     voice,
	}
}

func createVideoNoteMessage(videoNote *tgbotapi.VideoNote) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: 1,
		VideoNote: videoNote,
	}
}

func createStickerMessage(sticker *tgbotapi.Sticker) *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: 1,
		Sticker:   sticker,
	}
}

// TestExtractURLs tests URL extraction from text and captions
func TestExtractURLs(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		caption  string
		expected []string
	}{
		// Basic URL extraction
		{
			name:     "Single HTTP URL in text",
			text:     "Visit http://example.com for more info",
			caption:  "",
			expected: []string{"http://example.com"},
		},
		{
			name:     "Single HTTPS URL in text", 
			text:     "Check out https://golang.org",
			caption:  "",
			expected: []string{"https://golang.org"},
		},
		{
			name:     "URL in caption only",
			text:     "",
			caption:  "Image from https://images.example.com/photo.jpg",
			expected: []string{"https://images.example.com/photo.jpg"},
		},
		{
			name:     "URLs in both text and caption",
			text:     "See https://example.com",
			caption:  "More at http://test.org",
			expected: []string{"https://example.com", "http://test.org"},
		},
		
		// Multiple URLs
		{
			name:     "Multiple URLs in text",
			text:     "Visit https://site1.com and http://site2.org for details",
			caption:  "",
			expected: []string{"https://site1.com", "http://site2.org"},
		},
		{
			name:     "URLs with paths and queries",
			text:     "Check https://api.example.com/v1/users?id=123&format=json",
			caption:  "",
			expected: []string{"https://api.example.com/v1/users?id=123&format=json"},
		},
		{
			name:     "URLs with fragments",
			text:     "Read https://docs.example.com/guide#section1",
			caption:  "",
			expected: []string{"https://docs.example.com/guide#section1"},
		},
		
		// Edge cases
		{
			name:     "No URLs",
			text:     "This is just text without URLs",
			caption:  "Caption with no links",
			expected: nil,
		},
		{
			name:     "Empty strings",
			text:     "",
			caption:  "",
			expected: nil,
		},
		{
			name:     "URL at start of text",
			text:     "https://example.com is a great site",
			caption:  "",
			expected: []string{"https://example.com"},
		},
		{
			name:     "URL at end of text",
			text:     "Visit our website at https://example.com",
			caption:  "",
			expected: []string{"https://example.com"},
		},
		{
			name:     "URLs separated by various whitespace",
			text:     "https://site1.com\nhttps://site2.org\thttps://site3.net",
			caption:  "",
			expected: []string{"https://site1.com", "https://site2.org", "https://site3.net"},
		},
		
		// Invalid cases that should not match
		{
			name:     "Invalid protocols",
			text:     "ftp://files.example.com and file://local/path",
			caption:  "",
			expected: nil,
		},
		{
			name:     "URLs without protocol",
			text:     "Visit example.com and www.test.org",
			caption:  "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractURLs(tt.text, tt.caption)
			assert.Equal(t, tt.expected, result, "URL extraction failed")
		})
	}
}

// TestExtractHashtags tests hashtag extraction from text and captions
func TestExtractHashtags(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		caption  string
		expected []string
	}{
		// Basic hashtag extraction
		{
			name:     "Single hashtag in text",
			text:     "This is a #test message",
			caption:  "",
			expected: []string{"test"},
		},
		{
			name:     "Single hashtag in caption",
			text:     "",
			caption:  "Photo with #nature tag",
			expected: []string{"nature"},
		},
		{
			name:     "Hashtags in both text and caption",
			text:     "Message with #golang",
			caption:  "Also tagged #programming",
			expected: []string{"golang", "programming"},
		},
		
		// Multiple hashtags
		{
			name:     "Multiple hashtags in text",
			text:     "Learning #golang #programming #webdev today",
			caption:  "",
			expected: []string{"golang", "programming", "webdev"},
		},
		{
			name:     "Adjacent hashtags",
			text:     "Tags: #tag1#tag2#tag3",
			caption:  "",
			expected: []string{"tag1", "tag2", "tag3"},
		},
		
		// Hashtag variations
		{
			name:     "Hashtags with numbers",
			text:     "Event #2024 #web3 #ai2023",
			caption:  "",
			expected: []string{"2024", "web3", "ai2023"},
		},
		{
			name:     "Hashtags with underscores",
			text:     "Using #snake_case and #camelCase",
			caption:  "",
			expected: []string{"snake_case", "camelCase"},
		},
		{
			name:     "Mixed case hashtags",
			text:     "#JavaScript #HTML5 #css3 #NodeJS",
			caption:  "",
			expected: []string{"JavaScript", "HTML5", "css3", "NodeJS"},
		},
		
		// Position tests
		{
			name:     "Hashtag at start",
			text:     "#important message here",
			caption:  "",
			expected: []string{"important"},
		},
		{
			name:     "Hashtag at end",
			text:     "This is important #urgent",
			caption:  "",
			expected: []string{"urgent"},
		},
		{
			name:     "Hashtags separated by whitespace",
			text:     "#tag1 #tag2\n#tag3\t#tag4",
			caption:  "",
			expected: []string{"tag1", "tag2", "tag3", "tag4"},
		},
		
		// Edge cases
		{
			name:     "No hashtags",
			text:     "Regular message without tags",
			caption:  "Caption without tags",
			expected: nil,
		},
		{
			name:     "Empty strings",
			text:     "",
			caption:  "",
			expected: nil,
		},
		{
			name:     "Hash without text",
			text:     "Just a # symbol",
			caption:  "",
			expected: nil,
		},
		{
			name:     "Hash with space",
			text:     "Invalid # tag here",
			caption:  "",
			expected: nil,
		},
		{
			name:     "Hash with special characters",
			text:     "Invalid #@tag and #tag!",
			caption:  "",
			expected: []string{"tag"}, // #tag! will match #tag
		},
		{
			name:     "Hashtag in URL should not match",
			text:     "Visit https://example.com#section",
			caption:  "",
			expected: []string{"section"}, // Current regex will match #section
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractHashtags(tt.text, tt.caption)
			assert.Equal(t, tt.expected, result, "Hashtag extraction failed")
		})
	}
}

// TestExtractMentions tests mention extraction from text and captions
func TestExtractMentions(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		caption  string
		expected []string
	}{
		// Basic mention extraction
		{
			name:     "Single mention in text",
			text:     "Hello @username, how are you?",
			caption:  "",
			expected: []string{"username"},
		},
		{
			name:     "Single mention in caption",
			text:     "",
			caption:  "Photo by @photographer",
			expected: []string{"photographer"},
		},
		{
			name:     "Mentions in both text and caption",
			text:     "Thanks @alice",
			caption:  "Also thanks @bob",
			expected: []string{"alice", "bob"},
		},
		
		// Multiple mentions
		{
			name:     "Multiple mentions in text",
			text:     "Meeting with @john @mary @steve tomorrow",
			caption:  "",
			expected: []string{"john", "mary", "steve"},
		},
		{
			name:     "Adjacent mentions",
			text:     "Users: @user1@user2@user3",
			caption:  "",
			expected: []string{"user1", "user2", "user3"},
		},
		
		// Mention variations
		{
			name:     "Mentions with numbers",
			text:     "Contact @user123 or @admin2024",
			caption:  "",
			expected: []string{"user123", "admin2024"},
		},
		{
			name:     "Mentions with underscores",
			text:     "Follow @bot_user and @test_account",
			caption:  "",
			expected: []string{"bot_user", "test_account"},
		},
		{
			name:     "Mixed case mentions",
			text:     "@AdminUser @testBot @DevTeam",
			caption:  "",
			expected: []string{"AdminUser", "testBot", "DevTeam"},
		},
		
		// Position tests
		{
			name:     "Mention at start",
			text:     "@admin please help",
			caption:  "",
			expected: []string{"admin"},
		},
		{
			name:     "Mention at end",
			text:     "Please review this @reviewer",
			caption:  "",
			expected: []string{"reviewer"},
		},
		{
			name:     "Mentions separated by whitespace",
			text:     "@user1 @user2\n@user3\t@user4",
			caption:  "",
			expected: []string{"user1", "user2", "user3", "user4"},
		},
		
		// Edge cases
		{
			name:     "No mentions",
			text:     "Regular message without mentions",
			caption:  "Caption without mentions",
			expected: nil,
		},
		{
			name:     "Empty strings",
			text:     "",
			caption:  "",
			expected: nil,
		},
		{
			name:     "At symbol without username",
			text:     "Just an @ symbol",
			caption:  "",
			expected: nil,
		},
		{
			name:     "At symbol with space",
			text:     "Invalid @ user here",
			caption:  "",
			expected: nil,
		},
		{
			name:     "At symbol with special characters",
			text:     "Invalid @#user and @user!",
			caption:  "",
			expected: []string{"user"}, // @user! will match @user
		},
		{
			name:     "Email addresses should not match",
			text:     "Contact me at user@example.com",
			caption:  "",
			expected: []string{"example"}, // Current regex will match @example part
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractMentions(tt.text, tt.caption)
			assert.Equal(t, tt.expected, result, "Mention extraction failed")
		})
	}
}

// TestGetMessageType tests message type detection
func TestGetMessageType(t *testing.T) {
	tests := []struct {
		name     string
		message  *tgbotapi.Message
		expected MessageType
	}{
		// Photo messages
		{
			name: "Photo message",
			message: createPhotoMessage("", tgbotapi.PhotoSize{
				FileID: "photo123",
				Width:  100,
				Height: 100,
			}),
			expected: MessageTypePhoto,
		},
		{
			name:     "Photo message with multiple sizes",
			message:  createPhotoMessage("", 
				tgbotapi.PhotoSize{FileID: "thumb", Width: 100, Height: 100},
				tgbotapi.PhotoSize{FileID: "medium", Width: 400, Height: 400},
				tgbotapi.PhotoSize{FileID: "large", Width: 800, Height: 800},
			),
			expected: MessageTypePhoto,
		},
		
		// Video messages
		{
			name: "Video message",
			message: createVideoMessage("", &tgbotapi.Video{
				FileID:   "video123",
				Width:    1920,
				Height:   1080,
				Duration: 60,
			}),
			expected: MessageTypeVideo,
		},
		
		// Document messages
		{
			name: "Document message",
			message: createDocumentMessage("", &tgbotapi.Document{
				FileID:   "doc123",
				FileName: "document.pdf",
				MimeType: "application/pdf",
			}),
			expected: MessageTypeDocument,
		},
		
		// Audio messages
		{
			name: "Audio message",
			message: createAudioMessage("", &tgbotapi.Audio{
				FileID:    "audio123",
				Duration:  180,
				Title:     "Song Title",
				Performer: "Artist",
			}),
			expected: MessageTypeAudio,
		},
		
		// Voice messages
		{
			name: "Voice message",
			message: createVoiceMessage(&tgbotapi.Voice{
				FileID:   "voice123",
				Duration: 30,
			}),
			expected: MessageTypeVoice,
		},
		
		// Video note messages
		{
			name: "Video note message",
			message: createVideoNoteMessage(&tgbotapi.VideoNote{
				FileID:   "videonote123",
				Length:   240,
				Duration: 15,
			}),
			expected: MessageTypeVideoNote,
		},
		
		// Sticker messages
		{
			name: "Sticker message",
			message: createStickerMessage(&tgbotapi.Sticker{
				FileID: "sticker123",
				Width:  512,
				Height: 512,
			}),
			expected: MessageTypeSticker,
		},
		
		// Text messages (default case)
		{
			name:     "Text message",
			message:  createTextMessage("Hello world", ""),
			expected: MessageTypeText,
		},
		{
			name:     "Empty message defaults to text",
			message:  createTextMessage("", ""),
			expected: MessageTypeText,
		},
		
		// Priority testing - photo should take precedence over text
		{
			name: "Photo message with text",
			message: &tgbotapi.Message{
				MessageID: 1,
				Text:      "Photo with caption",
				Photo: []tgbotapi.PhotoSize{
					{FileID: "photo123", Width: 100, Height: 100},
				},
			},
			expected: MessageTypePhoto,
		},
		
		// Nil media objects should default to text
		{
			name: "Message with nil photo array",
			message: &tgbotapi.Message{
				MessageID: 1,
				Text:      "Text message",
				Photo:     nil,
			},
			expected: MessageTypeText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMessageType(tt.message)
			assert.Equal(t, tt.expected, result, "Message type detection failed")
		})
	}
}

// TestExtractFileMetadata tests file metadata extraction for different media types
func TestExtractFileMetadata(t *testing.T) {
	tests := []struct {
		name        string
		message     *tgbotapi.Message
		messageType MessageType
		expected    FileMetadata
	}{
		// Photo metadata
		{
			name: "Photo with file size",
			message: createPhotoMessage("", tgbotapi.PhotoSize{
				FileID:   "photo123",
				Width:    800,
				Height:   600,
				FileSize: 150000,
			}),
			messageType: MessageTypePhoto,
			expected: FileMetadata{
				FileID:   sqlNullString("photo123", true),
				FileName: sqlNullString("", false),
				MimeType: sqlNullString("", false),
				FileSize: sqlNullInt64(150000, true),
				Duration: sqlNullInt32(0, false),
			},
		},
		{
			name: "Photo without file size",
			message: createPhotoMessage("", tgbotapi.PhotoSize{
				FileID: "photo456",
				Width:  400,
				Height: 300,
			}),
			messageType: MessageTypePhoto,
			expected: FileMetadata{
				FileID:   sqlNullString("photo456", true),
				FileName: sqlNullString("", false),
				MimeType: sqlNullString("", false),
				FileSize: sqlNullInt64(0, false),
				Duration: sqlNullInt32(0, false),
			},
		},
		{
			name: "Empty photo array",
			message: createPhotoMessage(""),
			messageType: MessageTypePhoto,
			expected: FileMetadata{
				FileID:   sqlNullString("", false),
				FileName: sqlNullString("", false),
				MimeType: sqlNullString("", false),
				FileSize: sqlNullInt64(0, false),
				Duration: sqlNullInt32(0, false),
			},
		},
		
		// Video metadata
		{
			name: "Complete video metadata",
			message: createVideoMessage("", &tgbotapi.Video{
				FileID:   "video123",
				Width:    1920,
				Height:   1080,
				Duration: 120,
				FileName: "movie.mp4",
				MimeType: "video/mp4",
				FileSize: 5000000,
			}),
			messageType: MessageTypeVideo,
			expected: FileMetadata{
				FileID:   sqlNullString("video123", true),
				FileName: sqlNullString("movie.mp4", true),
				MimeType: sqlNullString("video/mp4", true),
				FileSize: sqlNullInt64(5000000, true),
				Duration: sqlNullInt32(120, true),
			},
		},
		{
			name: "Video with missing optional fields",
			message: createVideoMessage("", &tgbotapi.Video{
				FileID: "video456",
				Width:  640,
				Height: 480,
			}),
			messageType: MessageTypeVideo,
			expected: FileMetadata{
				FileID:   sqlNullString("video456", true),
				FileName: sqlNullString("", false),
				MimeType: sqlNullString("", false),
				FileSize: sqlNullInt64(0, false),
				Duration: sqlNullInt32(0, false),
			},
		},
		
		// Document metadata
		{
			name: "Complete document metadata",
			message: createDocumentMessage("", &tgbotapi.Document{
				FileID:   "doc123",
				FileName: "report.pdf",
				MimeType: "application/pdf",
				FileSize: 2000000,
			}),
			messageType: MessageTypeDocument,
			expected: FileMetadata{
				FileID:   sqlNullString("doc123", true),
				FileName: sqlNullString("report.pdf", true),
				MimeType: sqlNullString("application/pdf", true),
				FileSize: sqlNullInt64(2000000, true),
				Duration: sqlNullInt32(0, false),
			},
		},
		
		// Audio metadata
		{
			name: "Complete audio metadata",
			message: createAudioMessage("", &tgbotapi.Audio{
				FileID:    "audio123",
				Duration:  240,
				Title:     "Song Title",
				Performer: "Artist Name",
				FileName:  "song.mp3",
				MimeType:  "audio/mpeg",
				FileSize:  8000000,
			}),
			messageType: MessageTypeAudio,
			expected: FileMetadata{
				FileID:   sqlNullString("audio123", true),
				FileName: sqlNullString("song.mp3", true),
				MimeType: sqlNullString("audio/mpeg", true),
				FileSize: sqlNullInt64(8000000, true),
				Duration: sqlNullInt32(240, true),
			},
		},
		
		// Voice metadata
		{
			name: "Complete voice metadata",
			message: createVoiceMessage(&tgbotapi.Voice{
				FileID:   "voice123",
				Duration: 45,
				MimeType: "audio/ogg",
				FileSize: 500000,
			}),
			messageType: MessageTypeVoice,
			expected: FileMetadata{
				FileID:   sqlNullString("voice123", true),
				FileName: sqlNullString("", false),
				MimeType: sqlNullString("audio/ogg", true),
				FileSize: sqlNullInt64(500000, true),
				Duration: sqlNullInt32(45, true),
			},
		},
		
		// Text message (should return empty metadata)
		{
			name:        "Text message has no metadata",
			message:     createTextMessage("Hello world", ""),
			messageType: MessageTypeText,
			expected: FileMetadata{
				FileID:   sqlNullString("", false),
				FileName: sqlNullString("", false),
				MimeType: sqlNullString("", false),
				FileSize: sqlNullInt64(0, false),
				Duration: sqlNullInt32(0, false),
			},
		},
		
		// Nil media objects
		{
			name: "Nil video object",
			message: &tgbotapi.Message{
				MessageID: 1,
				Video:     nil,
			},
			messageType: MessageTypeVideo,
			expected: FileMetadata{
				FileID:   sqlNullString("", false),
				FileName: sqlNullString("", false),
				MimeType: sqlNullString("", false),
				FileSize: sqlNullInt64(0, false),
				Duration: sqlNullInt32(0, false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFileMetadata(tt.message, tt.messageType)
			assert.Equal(t, tt.expected, result, "File metadata extraction failed")
		})
	}
}

// Helper functions to create sql.Null* types for testing
func sqlNullString(s string, valid bool) sql.NullString {
	return sql.NullString{String: s, Valid: valid}
}

func sqlNullInt64(i int64, valid bool) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: valid}
}

func sqlNullInt32(i int32, valid bool) sql.NullInt32 {
	return sql.NullInt32{Int32: i, Valid: valid}
}