# Telegram Sticker & Emoji Stiller Bot

A Telegram bot written in Go that allows users to create copies of sticker packs and emoji packs under their ownership.

## Features

### Core Features
- 📦 **Copy Sticker Packs**: Create your own copy of any public sticker pack
- 😀 **Copy Emoji Packs**: Create your own copy of any public custom emoji pack
- 📊 **Pack Statistics**: View pack details including title and item count before creating
- 📋 **List Your Packs**: See all packs you've created with the bot
- 🗑️ **Delete Packs**: Remove packs from your list (via `/delete` command)
- 💾 **Persistent Storage**: All created packs are saved to a SQLite database
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

### Optional
- `PUBLIC_URL` - Your public URL for webhooks (e.g., `https://your-app.railway.app`)
  - If not set, uses long polling mode (for local development)
  - If set, uses webhook mode (for production)
- `PORT` - Server port for webhooks (default: `8443`, Railway sets this automatically)
- `ADMIN_IDS` - Comma-separated list of admin Telegram user IDs for broadcast feature
  - Example: `123456789,987654321`
  - Get your user ID from [@userinfobot](https://t.me/userinfobot)
- `DB_PATH` - Database file path (default: `./data/packs.db`)

## Development

### Prerequisites
- Go 1.25.0 or higher
- Telegram bot token from [@BotFather](https://t.me/BotFather)

### Local Setup (Polling Mode)

```powershell
# PowerShell (Windows)
$env:TOKEN="your_bot_token_here"
$env:ADMIN_IDS="your_telegram_user_id"

# Run the bot
go run main.go
```

```bash
# Bash (Linux/Mac)
export TOKEN="your_bot_token_here"
export ADMIN_IDS="your_telegram_user_id"

# Run the bot
go run main.go
```

The bot will automatically use **long polling mode** when `PUBLIC_URL` is not set.

### Project Structure

```
.
├── handlers/       # Request handlers
│   ├── pack.go       # Unified pack handler for stickers and emojis
│   └── admin.go      # Admin commands (broadcast, stats)
├── services/      # Business logic services
│   ├── download.go   # Download files from Telegram
│   ├── upload.go     # Upload and create sticker/emoji sets
│   ├── session.go    # Session management
│   └── telegram.go   # Telegram API interactions
├── db/            # Database layer
│   ├── models.go        # Data models (packs)
│   ├── user_tracking.go # User tracking model
│   ├── repository.go    # Database operations
│   └── schema.go        # Database schema
├── types/         # Type definitions
├── utils/         # Utility functions
├── i18n/          # Internationalization
│   ├── en.go         # English translations
│   └── ua.go         # Ukrainian translations
└── data/          # Data directory (temp files, database)
```

## Architecture

### Functional Programming Approach

This codebase follows functional programming patterns:
- Pure functions in utilities
- Service functions that take context as parameters
- No classes except for custom error types
- Heavy use of functional utilities

### Session Management

The bot uses an in-memory session store to track conversation state:
- `waiting_for_pack_name` - User has sent a pack link, waiting for new name

Session data includes:
- `OriginalItems` - Array of stickers/emojis from fetched pack
- `Title` - Original pack title
- `PackType` - Type of pack (sticker or emoji)

### Bot Flow

1. User sends sticker pack link (`t.me/addstickers/...`) or emoji pack link (`t.me/addemoji/...`)
2. Bot fetches pack details via Telegram API
3. Bot shows pack statistics (title, item count)
4. Bot asks for new pack name
5. User provides name (validated: non-empty, max 64 chars, alphanumeric + underscore)
6. Bot downloads all stickers/emojis to temp directory
7. Bot creates new sticker/emoji set using Telegram API
8. Bot saves pack info to database
9. Temp files are cleaned up
10. Session is cleared

## Deployment

### Railway (Production with Webhooks)

The bot is optimized for Railway deployment with webhooks:

1. **Create a new project on Railway**
2. **Add environment variables**:
   - `TOKEN` - Your bot token
   - `ADMIN_IDS` - Admin user IDs (comma-separated)
3. **Deploy from GitHub** (first deployment)
4. **Get your Railway app URL** from the dashboard
5. **Add the webhook URL**:
   - `RAILWAY_PUBLIC_DOMAIN` - Your Railway app URL (e.g., `https://your-app.railway.app`)
6. **Redeploy** - The bot will now use webhook mode

A volume will be automatically mounted at `/app/data` to persist the database.

📖 **Detailed guide**: See [docs/railway-webhook-setup.md](docs/railway-webhook-setup.md)

### Docker

Build and run with Docker:

```bash
docker build -t sticker-bot .
docker run -e TOKEN=your_token_here sticker-bot
```

### Technology Stack

- **Language**: Go 1.25.0
- **Framework**: Telebot v4 (`gopkg.in/telebot.v4`)
- **Database**: SQLite with optimized indexes
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
