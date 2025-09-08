# Telegram Organizer Mini-App - Frontend Development Instructions

## Project Context
You are working on a Telegram mini-app frontend for a content organizer bot. The project uses React and integrates with Telegram Web App SDK. The backend API is deployed on Yandex Cloud Functions with PostgreSQL database.

## Current Architecture
- **Backend**: Go Cloud Functions on Yandex Cloud
- **API Endpoint**: `/api/user/tags` (working, returns user tags with message counts)
- **Authentication**: Telegram Web App initData validation (working)
- **Database**: PostgreSQL with schema provided below

## Database Schema Reference
```sql
-- Users table
users (id, telegram_id, username, first_name, last_name, created_at, updated_at, is_active)

-- Messages table (core content)
messages (id, user_id, telegram_message_id, message_type, text_content, caption, 
         file_id, file_name, file_size, mime_type, duration, created_at, 
         forwarded_date, forwarded_from, urls, hashtags, mentions, search_vector)

-- Tags table
tags (id, user_id, name, color, created_at)

-- Message-Tags relationship
message_tags (id, message_id, tag_id, created_at)
```

## Current Frontend Structure
```
miniapp-frontend/
├── src/
│   ├── components/
│   │   ├── TagList.js      # Lists all user tags
│   │   ├── TagItem.js      # Individual tag display
│   │   └── Header.js       # App header
│   ├── App.js              # Main component
│   └── index.js           # Entry point
├── public/
└── package.json
```

## Required Improvements

### Phase 1. Error Handling & Debug Widgets
- **Health Check Widget**: Show only when there are API errors or connection issues
- **Debug Information Widget**: Display only during errors, show:
    - API response status
    - Network connectivity
    - Authentication status
    - Last error message with timestamp
- **Normal State**: Hide both widgets when everything works correctly

### Phase 2. Core Functionality Enhancement
**Tag Navigation Flow:**
```
TagList → Click Tag → MessageList (for that tag) → Click Message → Redirect to original Telegram message
```

**Required Components:**
- `MessageList.js` - Display messages for selected tag
- `MessageItem.js` - Individual message display with click-to-redirect
- Navigation state management between views

### Phase 3. API Endpoints to Implement
You need to create these new backend endpoints:

```go
// GET /api/user/tags/{tagId}/messages
// Returns messages for specific tag
type MessageResponse struct {
    ID                int64     `json:"id"`
    TelegramMessageID int64     `json:"telegram_message_id"`
    MessageType       string    `json:"message_type"`
    TextContent       *string   `json:"text_content"`
    Caption           *string   `json:"caption"`
    FileName          *string   `json:"file_name"`
    CreatedAt         time.Time `json:"created_at"`
    ForwardedFrom     *string   `json:"forwarded_from"`
    URLs              []string  `json:"urls"`
    Hashtags          []string  `json:"hashtags"`
}
```

### Phase 4. Navigation Implementation
**State Management:**
- Use React hooks (useState) or Context for navigation
- Track current view: 'tags' | 'messages' | 'loading' | 'error'
- Store selected tag information

**Telegram Message Redirection:**
```javascript
// Redirect to original message in Telegram
const redirectToMessage = (userId, messageId) => {
  const telegramUrl = `https://t.me/c/${userId}/${messageId}`;
  window.Telegram.WebApp.openTelegramLink(telegramUrl);
};
```

### Phase 5. UI/UX Requirements
- **Loading States**: Show spinners during API calls
- **Back Navigation**: Add back button when viewing messages
- **Empty States**: Handle cases with no messages for a tag
- **Error Recovery**: Allow retry buttons for failed requests
- **Responsive Design**: Work on mobile devices

### 6. Component Specifications

#### MessageList Component
```jsx
// Should display:
// - Tag name as header
// - Back button to return to tags
// - List of messages with preview
// - Message type icons
// - Creation date
// - Click handler for redirection
```

#### MessageItem Component
```jsx
// Should display:
// - Message preview (text/caption/filename)
// - Message type (text, photo, video, document, etc.)
// - Date/time
// - Forwarded from info (if applicable)
// - Click handler for Telegram redirection
```

#### Error Handling Components
```jsx
// HealthCheck Widget - only show on errors
// - Connection status
// - Last successful API call
// - Retry functionality

// Debug Widget - only show on errors  
// - Full error details
// - API endpoint called
// - Response status codes
// - Network timing info
```

## Success Criteria
1. ✅ Health/debug widgets only appear during errors
2. ✅ Clicking tag shows its messages
3. ✅ Clicking message redirects to original in Telegram
4. ✅ Smooth navigation between views
5. ✅ Proper error handling and recovery
6. ✅ Mobile-responsive design
7. ✅ Fast loading and good UX

## Existing Working Elements
- Telegram Web App SDK integration ✅
- Authentication with backend ✅
- Tag list display ✅
- Basic styling and layout ✅

Focus on creating a smooth, intuitive user experience where users can easily browse their organized content and quickly jump to original messages in Telegram.