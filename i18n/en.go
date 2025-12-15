package i18n

var En = map[string]string{
	"hello":   "Hello",
	"welcome": "Welcome to Sticker & Emoji Stiller @%s!\n\nCommands:\n/copy <link> - copy a pack\n/download <link> - download pack as ZIP\n/settings - set default action\n/language - change language\n\nOr just send:\n- Pack link (uses your default action)\n- Image (converts to sticker)\n- Sticker (converts to image)",
	"help":    "Commands:\n/copy <link> - copy a pack\n/download <link> - download pack as ZIP\n/settings - set default action\n/language - change language\n\nOr just send:\n- Pack link (uses your default action)\n- Image (converts to sticker)\n- Sticker (converts to image)",

	"start-command":    "Start (or restart) bot",
	"help-command":     "Show help message",
	"copy-command":     "Copy a sticker/emoji pack",
	"download-command": "Download pack as ZIP",
	"settings-command": "Set default action for links",

	"pack-stats":       "📦 Found %s pack: \"%s\"\n📊 Contains: %d items\n\nWhat would you like to name your new pack?\n\nType /cancel to cancel",
	"creating-pack":    "Creating your %s pack... This may take a while.",
	"downloading-pack": "📥 Downloading %s pack (%d items)...",
	"success":          "✅ Success! Your %s pack is ready:\n🔗 %s",
	"ask-pack-name":    "What would you like to name your %s pack? (Original: %s)\n\nJust type a name and I'll convert it to a valid format!\n\nType /cancel to cancel",
	"no-pack-data":     "No pack data found. Please start over.",
	"error":            "❌ Something went wrong. Please try again later.",
	"name-taken":       "This pack name is already taken. Please choose a different name or type /cancel to cancel.",

	"name-empty":         "Pack name cannot be empty. Please enter a valid name or type /cancel to cancel.",
	"name-too-long":      "Pack name is too long (max 64 characters). Please enter a shorter name or type /cancel to cancel.",
	"name-invalid-chars": "Pack name can only contain lowercase letters (a-z), numbers (0-9), and underscores (_). Please try again or type /cancel to cancel.",
	"cancelled":          "Operation cancelled.",

	"invalid-link": "Invalid link. Please send a valid sticker or emoji pack link, or use /copy or /download command.",
	"pack-type":    "sticker",
	"emoji-type":   "emoji",

	"cancel-command": "Cancel current operation",

	"settings-prompt":  "Current default action: %s\n\nSelect default action when you send a pack link:",
	"settings-saved":   "Settings saved!",
	"settings-updated": "✅ Default action set to: %s",
	"btn-copy":         "📋 Copy",
	"btn-download":     "📥 Download",

	"no-image":               "Please send an image.",
	"no-sticker":             "Please send a sticker.",
	"animated-not-supported": "Animated/video stickers cannot be converted to image.",
	"unsupported-format":     "Unsupported format. Please send a JPEG, PNG or WebP image.",

	"copy-usage":     "Usage: /copy <pack_link>\n\nExample: /copy https://t.me/addstickers/PackName",
	"download-usage": "Usage: /download <pack_link>\n\nExample: /download https://t.me/addstickers/PackName",

	"language-command": "Change interface language",
	"language-prompt":  "Current language: %s\n\nSelect your preferred language:",
	"language-saved":   "Language saved!",
	"language-updated": "✅ Language changed to: %s",
	"btn-lang-en":      "🇬🇧 English",
	"btn-lang-ua":      "🇺🇦 Українська",
	"btn-lang-pl":      "🇵🇱 Polski",
	"lang-name-en":     "English",
	"lang-name-ua":     "Ukrainian",
	"lang-name-pl":     "Polish",
}
