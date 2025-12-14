# Telegram Sticker & Emoji Stiller Bot

A Telegram bot written in Go that allows users to create copies of sticker packs and emoji packs under their ownership.

## Features

### Core Features
- 📦 **Copy Sticker Packs**: Create your own copy of any public sticker pack
- 😀 **Copy Emoji Packs**: Create your own copy of any public custom emoji pack
- 📊 **Pack Statistics**: View pack details including title and item count before creating
- 📋 **List Your Packs**: See all packs you've created with the bot
- 🗑️ **Delete Packs**: Remove packs from your list (via `/delete` command)
- 💾 **Persistent Storage**: All created packs are saved to a PostgreSQL database
- 🌍 **Multi-language**: Supports English and Ukrainian

## Commands

### Public Commands
- `/start` - Start or restart the bot
- `/help` - Show help message
- `/list` - List all packs you've created
- `/delete <pack_id>` - Delete a pack by its ID
- `/cancel` - Cancel current operation

### Admin Commands
- `/commands` - Force update bot commands
- `/stats` - View bot statistics

## Usage

1. Send the bot a sticker pack link (e.g., `t.me/addstickers/packname`) or emoji pack link (e.g., `t.me/addemoji/packname`)
2. The bot will show you pack statistics and ask for a name
3. Type a name for your new pack
4. Wait while the bot creates your pack
5. Receive the link to your new pack!

## Environment Variables

### Required
- `TOKEN` - Your Telegram bot token from [@BotFather](https://t.me/BotFather)
- `DATABASE_URL` - PostgreSQL connection string
  - Format: `postgresql://username:password@host:port/database?sslmode=require`
  - Railway provides this automatically when you connect a PostgreSQL database

### Optional
- `RAILWAY_PUBLIC_DOMAIN` - Your public URL for webhooks (e.g., `https://your-app.railway.app`)
  - If not set, uses long polling mode (for local development)
  - If set, uses webhook mode (for production)
- `RAILWAY_PORT` - Server port for webhooks (default: `8443`, Railway sets this automatically)
- `ADMIN_IDS` - Comma-separated list of admin Telegram user IDs for broadcast feature
  - Example: `123456789,987654321`
  - Get your user ID from [@userinfobot](https://t.me/userinfobot)

## Development

### Prerequisites
- Go 1.25.0 or higher
- PostgreSQL database
- Telegram bot token from [@BotFather](https://t.me/BotFather)

### Local Setup (Polling Mode)

```powershell
# PowerShell (Windows)
$env:TOKEN="your_bot_token_here"
$env:DATABASE_URL="postgresql://user:password@localhost:5432/dbname?sslmode=disable"
$env:ADMIN_IDS="your_telegram_user_id"

# Run the bot
go run main.go
```

```bash
# Bash (Linux/Mac)
export TOKEN="your_bot_token_here"
export DATABASE_URL="postgresql://user:password@localhost:5432/dbname?sslmode=disable"
export ADMIN_IDS="your_telegram_user_id"

# Run the bot
go run main.go
```

The bot will automatically use **long polling mode** when `RAILWAY_PUBLIC_DOMAIN` is not set.

## Deployment

### Railway (Production with Webhooks)

The bot is optimized for Railway deployment with webhooks:

1. **Create a PostgreSQL database on Railway**
2. **Create a new service for the bot**
3. **Connect the PostgreSQL database to your bot service**
   - Railway will automatically provide the `DATABASE_URL` environment variable
4. **Add environment variables**:
   - `TOKEN` - Your bot token
   - `ADMIN_IDS` - Admin user IDs (comma-separated)
5. **Deploy from GitHub** (first deployment)
6. **Get your Railway app URL** from the dashboard
7. **Add the webhook URL**:
   - `RAILWAY_PUBLIC_DOMAIN` - Your Railway app URL (e.g., `https://your-app.railway.app`)
8. **Redeploy** - The bot will now use webhook mode

📖 **Detailed guide**: See [docs/railway-webhook-setup.md](docs/railway-webhook-setup.md)

### Docker

Build and run with Docker:

```bash
docker build -t sticker-bot .
docker run -e TOKEN=your_token_here \
  -e DATABASE_URL=postgresql://user:password@host:5432/dbname \
  sticker-bot
```

### Technology Stack

- **Language**: Go 1.25.0
- **Framework**: Telebot v4 (`gopkg.in/telebot.v4`)
- **Database**: PostgreSQL (GORM with pgx driver)
- **Deployment**: Docker + Railway
- **Architecture**: Functional programming patterns

## Error Handling

- Users only see generic error messages to avoid confusion
- Detailed errors are logged for debugging
- All operations use retry logic for reliability

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Documentation

- [CHANGELOG.md](CHANGELOG.md) - Version history and changes

## License

MIT
