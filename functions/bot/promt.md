# Telegram Mini-App Implementation Metaprompt

You are tasked with implementing a basic Telegram Mini-App for the "Telegram Content Organizer" project. This mini-app should display a list of user's tags as the initial MVP feature.

## Project Context

**Architecture Overview:**
- Backend: Go with PostgreSQL database (already implemented)
- Bot: Telegram bot handles message processing and tagging (already implemented)
- Mini-App: React/Vue.js frontend hosted on Yandex Cloud Object Storage (to be implemented)
- API: Go-based API endpoints via Yandex Cloud Functions (to be implemented)

**Current Status:**
- ✅ Telegram bot is fully functional with tag management
- ✅ Database schema is implemented with full test coverage
- ❌ Mini-app frontend and API endpoints need to be created

## Technical Requirements

### 1. Database Schema (Already Implemented)
```sql
-- Users table
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

-- Tags table with user relationship
CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(telegram_id),
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7), -- hex color code
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

-- Message tags relationship (for future features)
CREATE TABLE message_tags (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT REFERENCES messages(id) ON DELETE CASCADE,
    tag_id BIGINT REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(message_id, tag_id)
);
```

### 2. Go Structs (Reference Implementation)
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

type Tag struct {
    ID        int64     `json:"id" db:"id"`
    UserID    int64     `json:"user_id" db:"user_id"`
    Name      string    `json:"name" db:"name"`
    Color     *string   `json:"color" db:"color"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

### 3. Database Connection
```
Connection String: take it from env variable DATABASE_URL  (check database.go file if needed).
```

## Implementation Requirements

### Phase 1: Backend API (Go Cloud Functions)

**File Structure:**
```
functions/
├── miniapp-api/
│   ├── main.go           # Lambda entry point
│   ├── handlers.go       # HTTP request handlers
│   ├── database.go       # DB operations (reuse from bot)
│   ├── auth.go           # Telegram Web App authentication
│   └── go.mod            # Dependencies
```

**Required API Endpoints:**

1. **GET /api/user/tags**
    - Purpose: Fetch all tags for authenticated user
    - Authentication: Telegram Web App `initData` validation
    - Response: JSON array of user's tags with message counts
    - SQL Query:
   ```sql
   SELECT t.*, COUNT(mt.message_id) as message_count
   FROM tags t
   LEFT JOIN message_tags mt ON t.id = mt.tag_id
   WHERE t.user_id = $1
   GROUP BY t.id
   ORDER BY message_count DESC;
   ```

**Authentication Method:**
- Validate Telegram Web App `initData` parameter
- Extract `user.id` from validated data
- Use as `user_id` for database queries

**Key Implementation Notes:**
- Use `github.com/gin-gonic/gin` for HTTP routing (consistent with project)
- Implement proper CORS headers for mini-app domain
- Use connection pooling for PostgreSQL
- Return proper HTTP status codes and error messages
- Implement request logging for debugging

### Phase 2: Frontend Mini-App (React/Vue.js)

**Technology Choice:** React (more suitable for Telegram Mini-Apps)

**File Structure:**
```
miniapp-frontend/
├── public/
│   ├── index.html
│   └── manifest.json
├── src/
│   ├── components/
│   │   ├── TagList.jsx    # Main tag list component
│   │   ├── TagItem.jsx    # Individual tag component
│   │   └── Header.jsx     # App header
│   ├── services/
│   │   └── api.js         # API communication
│   ├── utils/
│   │   └── telegram.js    # Telegram Web App utilities
│   ├── App.jsx            # Main app component
│   ├── index.js           # Entry point
│   └── styles.css         # Styling
├── package.json
└── README.md
```

**Core Components:**

1. **App.jsx** - Main application component
    - Initialize Telegram Web App SDK
    - Handle authentication state
    - Render TagList component

2. **TagList.jsx** - Display user's tags
    - Fetch tags from API on component mount
    - Show loading state during API calls
    - Display empty state if no tags exist
    - Show tag count for each tag

3. **TagItem.jsx** - Individual tag component
    - Display tag name with optional color indicator
    - Show message count
    - Responsive design for mobile

**Required Features:**
- Responsive design (mobile-first)
- Loading states and error handling
- Empty state when user has no tags
- Tag colors display (if set)
- Message count for each tag
- Telegram Web App theming integration

### Phase 3: Deployment Configuration

**Yandex Cloud Object Storage (Frontend):**
- Build optimized production bundle
- Upload to Object Storage bucket
- Configure as static website
- Set up HTTPS domain

**Yandex Cloud Functions (Backend API):**
- Package Go function with dependencies
- Configure environment variables
- Set up API Gateway routing
- Enable CORS for mini-app domain

## Specific Implementation Tasks

### Task 1: Telegram Web App Authentication
Implement secure authentication using Telegram's `initData` parameter:

```go
func validateTelegramWebApp(initData string, botToken string) (int64, error) {
    // Parse initData
    // Validate hash using bot token
    // Extract user information
    // Return user_id for database queries
}
```

### Task 2: Tags API Endpoint
Create endpoint that returns user's tags with message counts:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "work",
      "color": "#FF5733",
      "message_count": 15,
      "created_at": "2025-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "name": "personal",
      "color": null,
      "message_count": 8,
      "created_at": "2025-01-14T15:20:00Z"
    }
  ]
}
```

### Task 3: React Mini-App UI
Create clean, mobile-optimized interface:

**Design Requirements:**
- Follow Telegram Mini-App design guidelines
- Use Telegram's color scheme and theming
- Responsive grid layout for tags
- Loading spinners and error states
- Empty state with helpful message

**Example Tag Item Design:**
```
┌─────────────────────────────┐
│ ● work                  (15)│
│   Created: Jan 15, 2025     │
└─────────────────────────────┘
```

## Success Criteria

### MVP Completion Checklist:
- [ ] Backend API endpoint `/api/user/tags` working
- [ ] Telegram Web App authentication implemented
- [ ] Frontend displays user's tags with counts
- [ ] Responsive mobile design
- [ ] Error handling and loading states
- [ ] Deployed to Yandex Cloud (both frontend and backend)
- [ ] Mini-app accessible via Telegram bot command

### Testing Requirements:
- [ ] API endpoint returns correct data for test users
- [ ] Frontend handles empty state (no tags)
- [ ] Authentication works with real Telegram users
- [ ] Mobile responsiveness tested on different screen sizes
- [ ] Error scenarios handled gracefully

## Future Expansion Notes

This basic implementation sets foundation for:
- Message browsing by tag
- Tag editing and deletion
- Search functionality
- Export features
- Advanced tag management

Keep code modular and well-documented for easy expansion to these features.

## Technical Constraints

- **Budget**: Stay within Yandex Cloud free tier initially
- **Performance**: Optimize for mobile networks
- **Security**: Proper authentication and input validation
- **Scalability**: Design for 1K+ users from start
- **Maintainability**: Clean, testable code with documentation

---

**Expected Deliverables:**
1. Working Go API backend deployable to Yandex Cloud Functions
2. React frontend deployable to Yandex Cloud Object Storage
3. Documentation for deployment and configuration
4. Basic testing suite for API endpoints
5. Mini-app accessible through Telegram bot

**Estimated Development Time:** 2-3 days for experienced developer

**Priority:** This is the foundational feature that enables all future mini-app functionality.