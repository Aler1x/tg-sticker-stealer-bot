package handlers

import (
	"tg-sticker-stiller-bot/db"
	"tg-sticker-stiller-bot/utils"

	tg "gopkg.in/telebot.v4"
)

func HandleLanguage(ctx tg.Context, users *db.UserRepository) error {
	userID := ctx.Sender().ID

	user, err := users.GetByID(userID)
	if err != nil || user == nil {
		return ctx.Send(utils.T("en", "error"))
	}

	currentLang := user.LanguageCode
	if currentLang == "" {
		currentLang = "en"
	}

	langName := utils.T(currentLang, "lang-name-"+currentLang)

	keyboard := &tg.ReplyMarkup{}
	btnEn := keyboard.Data(utils.T(currentLang, "btn-lang-en"), "set_language", "en")
	btnUa := keyboard.Data(utils.T(currentLang, "btn-lang-ua"), "set_language", "ua")
	btnPl := keyboard.Data(utils.T(currentLang, "btn-lang-pl"), "set_language", "pl")

	keyboard.Inline(
		keyboard.Row(btnEn, btnUa, btnPl),
	)

	return ctx.Send(utils.T(currentLang, "language-prompt", langName), keyboard)
}

func HandleLanguageCallback(ctx tg.Context, users *db.UserRepository) error {
	userID := ctx.Sender().ID
	langCode := ctx.Data()

	if langCode != "en" && langCode != "ua" && langCode != "pl" {
		return ctx.Respond(&tg.CallbackResponse{Text: utils.T("en", "error")})
	}

	if err := users.SetLanguage(userID, langCode); err != nil {
		utils.Logger("error", "Failed to update language", map[string]any{
			"userId": userID,
			"error":  err.Error(),
		})
		return ctx.Respond(&tg.CallbackResponse{Text: utils.T(langCode, "error")})
	}

	langName := utils.T(langCode, "lang-name-"+langCode)
	ctx.Respond(&tg.CallbackResponse{Text: utils.T(langCode, "language-saved")})
	return ctx.Edit(utils.T(langCode, "language-updated", langName))
}
