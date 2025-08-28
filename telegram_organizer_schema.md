# Telegram Organizer - Database Schema

## Core Tables

### 1. Users
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);
```

### 2. Messages (Core Content)
```sql
CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(telegram_id),
    telegram_message_id BIGINT NOT NULL,
    message_type VARCHAR(50) NOT NULL, -- text, photo, video, document, audio, etc.
    text_content TEXT,
    caption TEXT,
    file_id VARCHAR(255), -- Telegram file_id for media
    file_name VARCHAR(255),
    file_size BIGINT,
    mime_type VARCHAR(100),
    duration INTEGER, -- for audio/video
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    forwarded_date TIMESTAMP,
    forwarded_from VARCHAR(255),
    
    -- Extracted metadata
    urls TEXT[],
    hashtags TEXT[],
    mentions TEXT[],
    
    -- Search optimization
    search_vector TSVECTOR,
    
    UNIQUE(user_id, telegram_message_id)
);
```

### 3. Tags
```sql
CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(telegram_id),
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7), -- hex color code
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, name)
);
```

### 4. Message Tags (Many-to-Many)
```sql
CREATE TABLE message_tags (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT REFERENCES messages(id) ON DELETE CASCADE,
    tag_id BIGINT REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(message_id, tag_id)
);
```

## Indexes
```sql
-- Search optimization
CREATE INDEX idx_messages_search_vector ON messages USING GIN(search_vector);
CREATE INDEX idx_messages_user_created ON messages(user_id, created_at DESC);
CREATE INDEX idx_messages_type ON messages(message_type);
CREATE INDEX idx_messages_hashtags ON messages USING GIN(hashtags);
CREATE INDEX idx_messages_urls ON messages USING GIN(urls);

-- Tag performance
CREATE INDEX idx_tags_user ON tags(user_id);
CREATE INDEX idx_message_tags_message ON message_tags(message_id);
CREATE INDEX idx_message_tags_tag ON message_tags(tag_id);

-- User lookups
CREATE INDEX idx_users_telegram_id ON users(telegram_id);
```

## Search Vector Functions & Triggers
```sql
-- Function to update search vector
CREATE OR REPLACE FUNCTION update_message_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := 
        setweight(to_tsvector('english', COALESCE(NEW.text_content, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.caption, '')), 'B') ||
        setweight(to_tsvector('english', array_to_string(NEW.hashtags, ' ')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update search vector
CREATE TRIGGER trigger_update_message_search_vector
    BEFORE INSERT OR UPDATE ON messages
    FOR EACH ROW
    EXECUTE FUNCTION update_message_search_vector();
```

## Connection String
```
postgres://tg_bot_user:your_password@66.248.207.105:5432/telegram_organizer?sslmode=disable
```

## Common Queries

### Search messages by text
```sql
SELECT * FROM messages 
WHERE user_id = $1 AND search_vector @@ plainto_tsquery($2)
ORDER BY ts_rank(search_vector, plainto_tsquery($2)) DESC, created_at DESC;
```

### Get messages by tag
```sql
SELECT m.* FROM messages m
JOIN message_tags mt ON m.id = mt.message_id
JOIN tags t ON mt.tag_id = t.id
WHERE m.user_id = $1 AND t.name = $2
ORDER BY m.created_at DESC;
```

### Get recent messages with tags
```sql
SELECT m.*, array_agg(t.name) as tag_names
FROM messages m
LEFT JOIN message_tags mt ON m.id = mt.message_id
LEFT JOIN tags t ON mt.tag_id = t.id
WHERE m.user_id = $1
GROUP BY m.id
ORDER BY m.created_at DESC
LIMIT $2;
```

### Get user tags with message counts
```sql
SELECT t.*, COUNT(mt.message_id) as message_count
FROM tags t
LEFT JOIN message_tags mt ON t.id = mt.tag_id
WHERE t.user_id = $1
GROUP BY t.id
ORDER BY message_count DESC;
```

## Go Struct Reference

### Basic structs for your backend:
```go
type User struct {
    ID          int64     `json:"id" db:"id"`
    TelegramID  int64     `json:"telegram_id" db:"telegram_id"`
    Username    *string   `json:"username" db:"username"`
    FirstName   *string   `json:"first_name" db:"first_name"`
    LastName    *string   `json:"last_name" db:"last_name"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    IsActive    bool      `json:"is_active" db:"is_active"`
}

type Message struct {
    ID                int64     `json:"id" db:"id"`
    UserID            int64     `json:"user_id" db:"user_id"`
    TelegramMessageID int64     `json:"telegram_message_id" db:"telegram_message_id"`
    MessageType       string    `json:"message_type" db:"message_type"`
    TextContent       *string   `json:"text_content" db:"text_content"`
    Caption           *string   `json:"caption" db:"caption"`
    FileID            *string   `json:"file_id" db:"file_id"`
    FileName          *string   `json:"file_name" db:"file_name"`
    FileSize          *int64    `json:"file_size" db:"file_size"`
    MimeType          *string   `json:"mime_type" db:"mime_type"`
    Duration          *int      `json:"duration" db:"duration"`
    CreatedAt         time.Time `json:"created_at" db:"created_at"`
    ForwardedDate     *time.Time `json:"forwarded_date" db:"forwarded_date"`
    ForwardedFrom     *string   `json:"forwarded_from" db:"forwarded_from"`
    URLs              []string  `json:"urls" db:"urls"`
    Hashtags          []string  `json:"hashtags" db:"hashtags"`
    Mentions          []string  `json:"mentions" db:"mentions"`
}

type Tag struct {
    ID        int64     `json:"id" db:"id"`
    UserID    int64     `json:"user_id" db:"user_id"`
    Name      string    `json:"name" db:"name"`
    Color     *string   `json:"color" db:"color"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type MessageTag struct {
    ID        int64     `json:"id" db:"id"`
    MessageID int64     `json:"message_id" db:"message_id"`
    TagID     int64     `json:"tag_id" db:"tag_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```