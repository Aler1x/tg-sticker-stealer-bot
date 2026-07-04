package handlers

import (
	"tg-sticker-stiller-bot/db"
	"tg-sticker-stiller-bot/services"
	"tg-sticker-stiller-bot/types"
	"tg-sticker-stiller-bot/utils"

	tg "gopkg.in/telebot.v4"
)

func flipPackType(packType types.StickerType) types.StickerType {
	if packType == types.StickerTypeEmoji {
		return types.StickerTypeRegular
	}
	return types.StickerTypeEmoji
}

func HandleTransformPack(ctx tg.Context, packName string, sourcePackType types.StickerType, bot *tg.Bot, sessions *services.SessionStore, users *db.UserRepository) error {
	userID := ctx.Sender().ID
	lang := utils.GetUserLanguage(users, userID, ctx.Message().Sender.LanguageCode)

	stickerSet, err := fetchStickerSet(bot, packName, sourcePackType, lang)
	if err != nil {
		return ctx.Send(utils.T(lang, "error"))
	}

	sourceTypeKey := "pack-type"
	if sourcePackType == types.StickerTypeEmoji {
		sourceTypeKey = "emoji-type"
	}

	targetType := flipPackType(sourcePackType)
	targetTypeKey := "pack-type"
	if targetType == types.StickerTypeEmoji {
		targetTypeKey = "emoji-type"
	}

	ctx.Send(utils.T(
		lang,
		"transform-stats",
		utils.T(lang, sourceTypeKey),
		stickerSet.Title,
		len(stickerSet.Stickers),
		utils.T(lang, targetTypeKey),
	))

	sessions.Set(userID, &services.Session{
		State:         services.StateWaitingForPackName,
		Action:        services.ActionTransform,
		Title:         stickerSet.Title,
		OriginalItems: stickerSet.Stickers,
		Name:          packName,
		PackType:      targetType,
	})

	return nil
}
