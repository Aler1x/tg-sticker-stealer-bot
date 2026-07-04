package i18n

var Pl = map[string]string{
	"hello":   "Cześć",
	"welcome": "Witaj w Sticker & Emoji Stiller @%s!\n\nKomendy:\n/copy <link> - skopiuj paczkę\n/transform <link> - konwertuj naklejki na emoji lub emoji na naklejki\n/download <link> - pobierz jako ZIP\n/settings - ustaw domyślną akcję\n/language - zmień język\n\nLub po prostu wyślij:\n- Link do paczki (użyje domyślnej akcji)\n- Obraz (konwertuje na naklejkę)\n- Naklejkę (konwertuje na obraz)",
	"help":    "Komendy:\n/copy <link> - skopiuj paczkę\n/transform <link> - konwertuj naklejki na emoji lub emoji na naklejki\n/download <link> - pobierz jako ZIP\n/settings - ustaw domyślną akcję\n/language - zmień język\n\nLub po prostu wyślij:\n- Link do paczki (użyje domyślnej akcji)\n- Obraz (konwertuje na naklejkę)\n- Naklejkę (konwertuje na obraz)",

	"start-command":    "Uruchom (lub zrestartuj) bota",
	"help-command":     "Pokaż wiadomość pomocy",
	"copy-command":     "Skopiuj paczkę naklejek/emoji",
	"transform-command": "Konwertuj naklejki na emoji lub emoji na naklejki",
	"download-command": "Pobierz paczkę jako ZIP",
	"settings-command": "Ustaw domyślną akcję dla linków",

	"pack-stats":       "📦 Znaleziono paczkę %s: \"%s\"\n📊 Zawiera: %d elementów\n\nJak chcesz nazwać swoją nową paczkę?\n\nWpisz /cancel aby anulować",
	"creating-pack":    "Tworzenie twojej paczki %s... To może chwilę potrwać.",
	"downloading-pack": "📥 Pobieranie paczki %s (%d elementów)...",
	"success":          "✅ Sukces! Twoja paczka %s jest gotowa:\n🔗 %s",
	"ask-pack-name":    "Jak chcesz nazwać swoją paczkę %s? (Oryginał: %s)\n\nPo prostu wpisz nazwę, a ja skonwertuję ją do poprawnego formatu!\n\nWpisz /cancel aby anulować",
	"no-pack-data":     "Nie znaleziono danych paczki. Zacznij od nowa.",
	"error":            "❌ Coś poszło nie tak. Spróbuj ponownie później.",
	"name-taken":       "Ta nazwa paczki jest już zajęta. Wybierz inną nazwę lub wpisz /cancel aby anulować.",

	"name-empty":         "Nazwa paczki nie może być pusta. Wprowadź poprawną nazwę lub wpisz /cancel aby anulować.",
	"name-too-long":      "Nazwa paczki jest za długa (maks. 64 znaki). Wprowadź krótszą nazwę lub wpisz /cancel aby anulować.",
	"name-invalid-chars": "Nazwa paczki może zawierać tylko małe litery (a-z), cyfry (0-9) i podkreślenia (_). Spróbuj ponownie lub wpisz /cancel aby anulować.",
	"cancelled":          "Operacja anulowana.",

	"invalid-link": "Nieprawidłowy link. Wyślij prawidłowy link do paczki naklejek lub emoji, lub użyj komendy /copy lub /download.",
	"pack-type":    "naklejek",
	"emoji-type":   "emoji",

	"cancel-command": "Anuluj bieżącą operację",

	"settings-prompt":  "Obecna domyślna akcja: %s\n\nWybierz domyślną akcję dla linków do paczek:",
	"settings-saved":   "Ustawienia zapisane!",
	"settings-updated": "✅ Domyślna akcja ustawiona na: %s",
	"btn-copy":         "📋 Kopiuj",
	"btn-download":     "📥 Pobierz",

	"no-image":               "Wyślij obraz.",
	"no-sticker":             "Wyślij naklejkę.",
	"animated-not-supported": "Animowane/wideo naklejki nie mogą być przekonwertowane na obraz.",
	"unsupported-format":     "Nieobsługiwany format. Wyślij obraz JPEG, PNG lub WebP.",

	"copy-usage":     "Użycie: /copy <link>\n\nPrzykład: /copy https://t.me/addstickers/PackName",
	"transform-usage": "Użycie: /transform <link>\n\nPrzykłady:\n/transform https://t.me/addstickers/PackName\n/transform https://t.me/addemoji/PackName",
	"transform-stats": "🔄 Znaleziono paczkę %s: \"%s\"\n📊 Zawiera: %d elementów\n\nZostanie przekonwertowana na paczkę %s.\n\nJak chcesz nazwać swoją nową paczkę?\n\nWpisz /cancel aby anulować",
	"download-usage": "Użycie: /download <link>\n\nPrzykład: /download https://t.me/addstickers/PackName",

	"language-command": "Zmień język interfejsu",
	"language-prompt":  "Obecny język: %s\n\nWybierz preferowany język:",
	"language-saved":   "Język zapisany!",
	"language-updated": "✅ Język zmieniony na: %s",
	"btn-lang-en":      "🇬🇧 English",
	"btn-lang-ua":      "🇺🇦 Українська",
	"btn-lang-pl":      "🇵🇱 Polski",
	"lang-name-en":     "English",
	"lang-name-ua":     "Українська",
	"lang-name-pl":     "Polski",
}
