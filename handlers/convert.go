package handlers

import (
	"os"
	"tg-sticker-stiller-bot/services"
	"tg-sticker-stiller-bot/utils"

	tg "gopkg.in/telebot.v4"
)

func HandleImageToSticker(ctx tg.Context, bot *tg.Bot) error {
	lang := ctx.Message().Sender.LanguageCode
	photo := ctx.Message().Photo

	if photo == nil {
		return ctx.Send(utils.T(lang, "no-image"))
	}

	stickerPath, err := services.ConvertImageToSticker(bot, photo)
	if err != nil {
		utils.Logger("error", "Failed to convert image to sticker", map[string]any{
			"userId": ctx.Sender().ID,
			"error":  err.Error(),
		})
		return ctx.Send(utils.T(lang, "error"))
	}
	defer os.Remove(stickerPath)

	sticker := &tg.Sticker{
		File: tg.FromDisk(stickerPath),
	}

	return ctx.Send(sticker)
}

func HandleStickerToImage(ctx tg.Context, bot *tg.Bot) error {
	lang := ctx.Message().Sender.LanguageCode
	sticker := ctx.Message().Sticker

	if sticker == nil {
		return ctx.Send(utils.T(lang, "no-sticker"))
	}

	if sticker.Animated || sticker.Video {
		return ctx.Send(utils.T(lang, "animated-not-supported"))
	}

	imagePath, err := services.ConvertStickerToImage(bot, sticker)
	if err != nil {
		utils.Logger("error", "Failed to convert sticker to image", map[string]any{
			"userId": ctx.Sender().ID,
			"error":  err.Error(),
		})
		return ctx.Send(utils.T(lang, "error"))
	}
	defer os.Remove(imagePath)

	doc := &tg.Document{
		File:     tg.FromDisk(imagePath),
		FileName: "sticker.png",
	}

	return ctx.Send(doc)
}
