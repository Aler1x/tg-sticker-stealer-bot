package handlers

import (
	"fmt"
	"os"
	"tg-sticker-stiller-bot/db"
	"tg-sticker-stiller-bot/services"
	"tg-sticker-stiller-bot/types"
	"tg-sticker-stiller-bot/utils"

	tg "gopkg.in/telebot.v4"
)

func fetchStickerSet(bot *tg.Bot, packName string, packType types.StickerType, lang string) (*types.StickerSet, error) {
	if packType == types.StickerTypeEmoji {
		emojiSet, err := services.FetchEmojiSet(bot, packName)
		if err != nil {
			utils.Logger("error", "Error fetching emoji pack", map[string]any{
				"packName": packName,
				"error":    err.Error(),
			})
			return nil, err
		}
		return &types.StickerSet{
			Name:     emojiSet.Name,
			Title:    emojiSet.Title,
			Stickers: emojiSet.Stickers,
		}, nil
	}

	stickerSet, err := services.FetchStickerSet(bot, packName)
	if err != nil {
		utils.Logger("error", "Error fetching sticker pack", map[string]any{
			"packName": packName,
			"error":    err.Error(),
		})
		return nil, err
	}
	return stickerSet, nil
}

func HandleCopyPack(ctx tg.Context, packName string, packType types.StickerType, bot *tg.Bot, sessions *services.SessionStore, users *db.UserRepository) error {
	userID := ctx.Sender().ID
	lang := utils.GetUserLanguage(users, userID, ctx.Message().Sender.LanguageCode)

	stickerSet, err := fetchStickerSet(bot, packName, packType, lang)
	if err != nil {
		return ctx.Send(utils.T(lang, "error"))
	}

	packTypeKey := "pack-type"
	if packType == types.StickerTypeEmoji {
		packTypeKey = "emoji-type"
	}

	ctx.Send(utils.T(lang, "pack-stats", utils.T(lang, packTypeKey), stickerSet.Title, len(stickerSet.Stickers)))

	sessions.Set(userID, &services.Session{
		State:         services.StateWaitingForPackName,
		Action:        services.ActionCopy,
		Title:         stickerSet.Title,
		OriginalItems: stickerSet.Stickers,
		Name:          packName,
		PackType:      packType,
	})

	return nil
}

func HandleDownloadPack(ctx tg.Context, packName string, packType types.StickerType, bot *tg.Bot, users *db.UserRepository) error {
	userID := ctx.Sender().ID
	lang := utils.GetUserLanguage(users, userID, ctx.Message().Sender.LanguageCode)

	stickerSet, err := fetchStickerSet(bot, packName, packType, lang)
	if err != nil {
		return ctx.Send(utils.T(lang, "error"))
	}

	packTypeKey := "pack-type"
	if packType == types.StickerTypeEmoji {
		packTypeKey = "emoji-type"
	}

	progressMsg, err := bot.Send(ctx.Recipient(), utils.T(lang, "downloading-pack", utils.T(lang, packTypeKey), len(stickerSet.Stickers)))
	if err != nil {
		utils.Logger("warn", "Failed to send progress message", map[string]any{"error": err.Error()})
	}

	progressCallback := func(current, total int) {
		if progressMsg != nil {
			newText := fmt.Sprintf("📥 Downloading: %d/%d items...", current, total)
			bot.Edit(progressMsg, newText)
		}
	}

	zipPath, err := services.CreateStickerZip(bot, stickerSet.Stickers, stickerSet.Name, progressCallback)
	if err != nil {
		if progressMsg != nil {
			bot.Delete(progressMsg)
		}
		utils.Logger("error", "Failed to create zip", map[string]any{"error": err.Error()})
		return ctx.Send(utils.T(lang, "error"))
	}
	defer os.Remove(zipPath)

	if progressMsg != nil {
		bot.Delete(progressMsg)
	}

	doc := &tg.Document{
		File:     tg.FromDisk(zipPath),
		FileName: fmt.Sprintf("%s.zip", stickerSet.Name),
	}

	return ctx.Send(doc)
}

func HandlePackNameInput(ctx tg.Context, userInput string, bot *tg.Bot, sessions *services.SessionStore, database *db.DB) error {
	userID := ctx.Sender().ID
	lang := utils.GetUserLanguage(database.Users, userID, ctx.Message().Sender.LanguageCode)

	session := sessions.Get(userID)

	if len(session.OriginalItems) == 0 {
		sessions.Clear(userID)
		return ctx.Send(utils.T(lang, "no-pack-data"))
	}

	normalizedName := utils.NormalizePackName(userInput)

	if !utils.ValidateNormalizedName(normalizedName) {
		errKey := utils.GetValidationError(normalizedName)
		return ctx.Send(utils.T(lang, errKey))
	}

	packTypeKey := "pack-type"
	if session.PackType == types.StickerTypeEmoji {
		packTypeKey = "emoji-type"
	}

	progressMsg, err := ctx.Bot().Send(ctx.Recipient(), utils.T(lang, "creating-pack", utils.T(lang, packTypeKey)))
	if err != nil {
		utils.Logger("warn", "Failed to send progress message", map[string]any{"error": err.Error()})
	}

	progressCallback := func(current, total int) {
		if progressMsg != nil {
			newText := fmt.Sprintf("📦 Processing: %d/%d items...", current, total)
			_, err := ctx.Bot().Edit(progressMsg, newText)
			if err != nil {
				utils.Logger("debug", "Failed to update progress", map[string]any{"error": err.Error()})
			}
		}
	}

	packLink, err := services.CreateStickerSet(bot, userID, bot.Me.Username, userInput, session.OriginalItems, session.PackType, progressCallback)
	if err != nil {
		if progressMsg != nil {
			ctx.Bot().Delete(progressMsg)
		}
		if botErr, ok := err.(*utils.BotError); ok {
			if botErr.I18nKey == "name-taken" {
				return ctx.Send(utils.T(lang, "name-taken"))
			}
		}
		sessions.Clear(userID)
		utils.Logger("error", "Error creating sticker set", map[string]any{
			"userId": userID,
			"error":  err.Error(),
		})
		return ctx.Send(utils.T(lang, "error"))
	}

	if progressMsg != nil {
		ctx.Bot().Delete(progressMsg)
	}

	if err := database.PackCreations.Record(userID, packLink, string(session.PackType)); err != nil {
		utils.Logger("error", "Failed to record pack creation analytics", map[string]any{
			"userId": userID,
			"error":  err.Error(),
		})
	}

	ctx.Send(utils.T(lang, "success", utils.T(lang, packTypeKey), packLink))
	sessions.Clear(userID)
	return nil
}


