
# Telegram Content Organizer - Project Brief

## Core Concept

A Telegram bot + mini-app that helps users organize their Saved Messages with tags, search, and smart categorization. Users forward messages to the bot for organization, while original content stays in Telegram.

## Architecture
```
Yandex Cloud (FREE tier):
├── Cloud Function: Bot Backend (Go)
├── Cloud Function: Mini-App API (Go) 
├── Object Storage: Mini-App Frontend (React/Vue)
└── API Gateway: HTTPS endpoints
Your VPS ($5-10/month):
└── PostgreSQL Database
```

## Tech Stack


- Backend: Go + Gin framework
- Bot: go-telegram-bot-api/telegram-bot-api/v5
- Database: PostgreSQL on VPS
- Frontend: React/Vue.js for mini-app
- Hosting: Yandex Cloud Functions + Object Storage
- Security: Environment variables → Yandex Lockbox

# Core Features

## MVP:

- Receive forwarded messages
- Extract metadata (text, type, URLs, hashtags)
- Store with manual tags
- Basic search by tags/date
- Simple mini-app dashboard

## Scaling Estimates

- 1K users: FREE (Yandex Cloud free tier)
- 10K users: ~$10-30/month
- 100K users: ~$50-150/month

## Security

- Environment variables for non-sensitive config
- Yandex Lockbox for secrets (DB password, bot token)
- Dedicated PostgreSQL user with minimal permissions
- SSL connections to database

# Key Implementation Notes

- Store only metadata, not actual message content
- Use webhooks (not polling) for bot
- Connection pooling for database
- User workflow: Forward to bot → Auto-analyze → Store metadata → Browse via mini-app
