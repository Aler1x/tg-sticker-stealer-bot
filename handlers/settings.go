package handlers

import (
	"tg-sticker-stiller-bot/db"
	"tg-sticker-stiller-bot/utils"

	tg "gopkg.in/telebot.v4"
)

func HandleSettings(ctx tg.Context, users *db.UserRepository) error {
	userID := ctx.Sender().ID
	lang := utils.GetUserLanguage(users, userID, ctx.Message().Sender.LanguageCode)

	user, err := users.GetByID(userID)
	if err != nil || user == nil {
		return ctx.Send(utils.T(lang, "error"))
	}

	currentAction := user.DefaultAction
	if currentAction == "" {
		currentAction = db.DefaultActionCopy
	}

	keyboard := &tg.ReplyMarkup{}
	btnCopy := keyboard.Data(utils.T(lang, "btn-copy"), "set_action", "copy")
	btnDownload := keyboard.Data(utils.T(lang, "btn-download"), "set_action", "download")

	keyboard.Inline(
		keyboard.Row(btnCopy, btnDownload),
	)

	return ctx.Send(utils.T(lang, "settings-prompt", currentAction), keyboard)
}

func HandleSettingsCallback(ctx tg.Context, users *db.UserRepository) error {
	userID := ctx.Sender().ID
	lang := utils.GetUserLanguage(users, userID, ctx.Sender().LanguageCode)
	action := ctx.Data()

	var dbAction db.DefaultAction
	switch action {
	case "copy":
		dbAction = db.DefaultActionCopy
	case "download":
		dbAction = db.DefaultActionDownload
	default:
		return ctx.Respond(&tg.CallbackResponse{Text: utils.T(lang, "error")})
	}

	if err := users.SetDefaultAction(userID, dbAction); err != nil {
		utils.Logger("error", "Failed to update default action", map[string]any{
			"userId": userID,
			"error":  err.Error(),
		})
		return ctx.Respond(&tg.CallbackResponse{Text: utils.T(lang, "error")})
	}

	ctx.Respond(&tg.CallbackResponse{Text: utils.T(lang, "settings-saved")})
	return ctx.Edit(utils.T(lang, "settings-updated", action))
}
