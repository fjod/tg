# Telegram Content Organizer Bot

A webhook-based Telegram bot for testing integration with the Telegram API.

## Features

- Webhook-based (production-ready architecture)
- Responds to `/start` and `/help` commands
- Echoes back received text messages
- Acknowledges forwarded messages
- HTTP health check endpoint
- Basic error handling and logging

## Setup

1. Create a new bot with [@BotFather](https://t.me/botfather) on Telegram
2. Copy the bot token
3. Create a `.env` file from the example:
   ```bash
   cp .env.example .env
   ```
4. Edit `.env` and configure your settings:
   ```
   TELEGRAM_BOT_TOKEN=your_actual_bot_token_here
   WEBHOOK_URL=https://your-domain.com
   PORT=8080
   ```

## Deployment Options

### Local Testing (with ngrok)
```bash
# Install ngrok and expose local port
ngrok http 8080

# Copy the https URL (e.g., https://abc123.ngrok.io) to WEBHOOK_URL in .env
# Then run the bot
go run main.go
```

### Production Deployment
```bash
# Set WEBHOOK_URL to your production domain
# Deploy to your server/cloud function
go run main.go
```

## Testing

1. Start the bot (ensure webhook is accessible)
2. Send `/start` command to your bot
3. Send any text message - bot will echo it back
4. Forward any message to the bot - it will acknowledge receipt
5. Check health endpoint: `GET /health`

## Environment Variables

- `TELEGRAM_BOT_TOKEN` - Your bot token from @BotFather (required)
- `WEBHOOK_URL` - Your public HTTPS URL (required)
- `PORT` - Server port (default: 8080)

## Notes

- **Webhook-based**: More efficient than polling, suitable for production
- **HTTPS required**: Telegram requires HTTPS for webhooks
- **Health check**: `/health` endpoint for monitoring
- **The `.env` file is automatically loaded if present**
- **Falls back to system environment variables if no `.env` file found**
- **The `.env` file is excluded from git to keep your secrets secure**

## Architecture Benefits

- ✅ **Production-ready**: Aligns with Yandex Cloud Functions architecture from README.md
- ✅ **Efficient**: No polling, Telegram pushes updates directly
- ✅ **Scalable**: Stateless HTTP server, easy to deploy anywhere
- ✅ **Monitoring**: Built-in health check endpoint