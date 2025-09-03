package main

import (
	"database/sql"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageType represents the type of a Telegram message
type MessageType string

const (
	MessageTypeText      MessageType = "text"
	MessageTypePhoto     MessageType = "photo"
	MessageTypeVideo     MessageType = "video"
	MessageTypeDocument  MessageType = "document"
	MessageTypeAudio     MessageType = "audio"
	MessageTypeVoice     MessageType = "voice"
	MessageTypeVideoNote MessageType = "video_note"
	MessageTypeSticker   MessageType = "sticker"
)

// FileMetadata contains file information extracted from a Telegram message
type FileMetadata struct {
	FileID   sql.NullString
	FileName sql.NullString
	MimeType sql.NullString
	FileSize sql.NullInt64
	Duration sql.NullInt32
}

func extractURLs(text, caption string) []string {
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	var urls []string
	if text != "" {
		urls = append(urls, urlRegex.FindAllString(text, -1)...)
	}
	if caption != "" {
		urls = append(urls, urlRegex.FindAllString(caption, -1)...)
	}
	return urls
}

func extractHashtags(text, caption string) []string {
	hashtagRegex := regexp.MustCompile(`#\w+`)
	var hashtags []string
	if text != "" {
		hashtags = append(hashtags, hashtagRegex.FindAllString(text, -1)...)
	}
	if caption != "" {
		hashtags = append(hashtags, hashtagRegex.FindAllString(caption, -1)...)
	}
	for i, tag := range hashtags {
		hashtags[i] = strings.TrimPrefix(tag, "#")
	}
	return hashtags
}

func extractMentions(text, caption string) []string {
	mentionRegex := regexp.MustCompile(`@\w+`)
	var mentions []string
	if text != "" {
		mentions = append(mentions, mentionRegex.FindAllString(text, -1)...)
	}
	if caption != "" {
		mentions = append(mentions, mentionRegex.FindAllString(caption, -1)...)
	}
	for i, mention := range mentions {
		mentions[i] = strings.TrimPrefix(mention, "@")
	}
	return mentions
}

func getMessageType(message *tgbotapi.Message) MessageType {
	if message.Photo != nil {
		return MessageTypePhoto
	}
	if message.Video != nil {
		return MessageTypeVideo
	}
	if message.Document != nil {
		return MessageTypeDocument
	}
	if message.Audio != nil {
		return MessageTypeAudio
	}
	if message.Voice != nil {
		return MessageTypeVoice
	}
	if message.VideoNote != nil {
		return MessageTypeVideoNote
	}
	if message.Sticker != nil {
		return MessageTypeSticker
	}
	return MessageTypeText
}

func extractFileMetadata(message *tgbotapi.Message, messageType MessageType) FileMetadata {
	var metadata FileMetadata

	switch messageType {
	case MessageTypePhoto:
		if len(message.Photo) > 0 {
			photo := message.Photo[0] // Get the smallest photo (thumbnail)
			metadata.FileID = sql.NullString{String: photo.FileID, Valid: true}
			if photo.FileSize != 0 {
				metadata.FileSize = sql.NullInt64{Int64: int64(photo.FileSize), Valid: true}
			}
		}
	case MessageTypeVideo:
		if message.Video != nil {
			metadata.FileID = sql.NullString{String: message.Video.FileID, Valid: true}
			if message.Video.FileName != "" {
				metadata.FileName = sql.NullString{String: message.Video.FileName, Valid: true}
			}
			if message.Video.MimeType != "" {
				metadata.MimeType = sql.NullString{String: message.Video.MimeType, Valid: true}
			}
			if message.Video.FileSize != 0 {
				metadata.FileSize = sql.NullInt64{Int64: int64(message.Video.FileSize), Valid: true}
			}
			if message.Video.Duration != 0 {
				metadata.Duration = sql.NullInt32{Int32: int32(message.Video.Duration), Valid: true}
			}
		}
	case MessageTypeDocument:
		if message.Document != nil {
			metadata.FileID = sql.NullString{String: message.Document.FileID, Valid: true}
			if message.Document.FileName != "" {
				metadata.FileName = sql.NullString{String: message.Document.FileName, Valid: true}
			}
			if message.Document.MimeType != "" {
				metadata.MimeType = sql.NullString{String: message.Document.MimeType, Valid: true}
			}
			if message.Document.FileSize != 0 {
				metadata.FileSize = sql.NullInt64{Int64: int64(message.Document.FileSize), Valid: true}
			}
		}
	case MessageTypeAudio:
		if message.Audio != nil {
			metadata.FileID = sql.NullString{String: message.Audio.FileID, Valid: true}
			if message.Audio.FileName != "" {
				metadata.FileName = sql.NullString{String: message.Audio.FileName, Valid: true}
			}
			if message.Audio.MimeType != "" {
				metadata.MimeType = sql.NullString{String: message.Audio.MimeType, Valid: true}
			}
			if message.Audio.FileSize != 0 {
				metadata.FileSize = sql.NullInt64{Int64: int64(message.Audio.FileSize), Valid: true}
			}
			if message.Audio.Duration != 0 {
				metadata.Duration = sql.NullInt32{Int32: int32(message.Audio.Duration), Valid: true}
			}
		}
	case MessageTypeVoice:
		if message.Voice != nil {
			metadata.FileID = sql.NullString{String: message.Voice.FileID, Valid: true}
			if message.Voice.MimeType != "" {
				metadata.MimeType = sql.NullString{String: message.Voice.MimeType, Valid: true}
			}
			if message.Voice.FileSize != 0 {
				metadata.FileSize = sql.NullInt64{Int64: int64(message.Voice.FileSize), Valid: true}
			}
			if message.Voice.Duration != 0 {
				metadata.Duration = sql.NullInt32{Int32: int32(message.Voice.Duration), Valid: true}
			}
		}
	}

	return metadata
}
