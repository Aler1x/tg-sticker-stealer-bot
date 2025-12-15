package utils

import (
	"fmt"
	"tg-sticker-stiller-bot/i18n"
)

func T(lang string, key string, args ...interface{}) string {
	template := getTemplate(lang, key)

	if len(args) == 0 {
		return template
	}

	return fmt.Sprintf(template, args...)
}

func getTemplate(lang string, key string) string {
	switch lang {
	case "ua":
		if val, ok := i18n.Ua[key]; ok {
			return val
		}
		return i18n.En[key]
	case "pl":
		if val, ok := i18n.Pl[key]; ok {
			return val
		}
		return i18n.En[key]
	case "en":
		return i18n.En[key]
	default:
		return i18n.En[key]
	}
}
