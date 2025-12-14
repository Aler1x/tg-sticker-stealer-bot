package utils

import "tg-sticker-stiller-bot/db"

func GetUserLanguage(users *db.UserRepository, userID int64, fallbackLang string) string {
	user, err := users.GetByID(userID)
	if err != nil || user == nil {
		return fallbackLang
	}

	if user.LanguageCode != "" {
		return user.LanguageCode
	}

	return fallbackLang
}
