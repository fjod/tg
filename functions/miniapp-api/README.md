# Telegram Mini-App API

This is the backend API service for the Telegram Mini-App, implementing Phase 1 of the project requirements.

## Features

- **GET /api/user/tags** - Fetch user's tags with message counts
- **Telegram Web App Authentication** - Secure validation using initData
- **CORS Support** - Ready for frontend integration
- **Lambda Compatible** - Deployable to Yandex Cloud Functions

## File Structure

```
miniapp-api/
├── main.go           # Lambda entry point + HTTP routing
├── handlers.go       # API endpoint handlers
├── database.go       # Database operations and structs
├── auth.go           # Telegram Web App authentication
├── main_test.go      # Basic tests
├── go.mod            # Dependencies
└── README.md         # This file
```

## API Endpoints

### GET /api/user/tags

Returns user's tags with message counts, sorted by message count (descending).

**Headers:**
- `Authorization: Bearer <telegram_initData>`

**Response Format:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 123456,
      "name": "work",
      "color": "#FF5733",
      "message_count": 15,
      "created_at": "2025-01-15T10:30:00Z"
    }
  ]
}
```

## Authentication

Uses Telegram Web App `initData` validation:
1. Parse initData parameters
2. Validate HMAC signature using bot token
3. Extract user ID for database queries

## Database Schema

Reuses existing schema from bot implementation:
- `users` table for user information
- `tags` table for user tags
- `message_tags` for tag-message relationships

## Testing

```bash
cd miniapp-api
go test -v
```

Note: Tests require `DATABASE_URL` environment variable for full coverage.

## Deployment

This service is designed for deployment to Yandex Cloud Functions with:
- Environment variables: `DATABASE_URL`, `TELEGRAM_BOT_TOKEN`
- Runtime: Go 1.23+
- Handler: `main.Handler`

## Next Steps

- Deploy to Yandex Cloud Functions
- Configure API Gateway routing
- Implement frontend (Phase 2)
- Add integration tests with real database