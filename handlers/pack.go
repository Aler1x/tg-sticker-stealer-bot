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

func HandlePackNameInput(ctx tg.Context, userInput string, bot *tg.Bot, sessions *services.SessionStore, packs *db.PackRepository, users *db.UserRepository) error {
	userID := ctx.Sender().ID
	lang := utils.GetUserLanguage(users, userID, ctx.Message().Sender.LanguageCode)

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

	packLink, err := services.CreateStickerSet(bot, userID, bot.Me.Username, userInput, session.OriginalItems, session.PackType, packs, progressCallback)
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

	ctx.Send(utils.T(lang, "success", utils.T(lang, packTypeKey), packLink))
	sessions.Clear(userID)
	return nil
}

func HandleListPacks(ctx tg.Context, page int, packs *db.PackRepository, users *db.UserRepository) error {
	userID := ctx.Sender().ID
	lang := utils.GetUserLanguage(users, userID, ctx.Message().Sender.LanguageCode)

	const pageSize = 5

	packList, total, err := packs.GetPaginated(userID, page, pageSize)
	if err != nil {
		utils.Logger("error", "Error getting packs for user", map[string]any{
			"userId": userID,
			"error":  err.Error(),
		})
		return ctx.Send(utils.T(lang, "error"))
	}

	if len(packList) == 0 {
		if page == 1 {
			return ctx.Send(utils.T(lang, "list-empty"))
		}
		return ctx.Send(utils.T(lang, "list-page-empty"))
	}

	totalPages := (total + pageSize - 1) / pageSize
	message := utils.T(lang, "list-header-paginated", page, totalPages, total)

	startIndex := (page - 1) * pageSize
	for i, pack := range packList {
		orderNum := startIndex + i + 1
		message += utils.T(lang, "list-item", orderNum, pack.PackTitle, pack.PackType, pack.StickerCount, pack.PackLink)
	}

	if totalPages > 1 {
		keyboard := buildPaginationKeyboard(page, totalPages)
		return ctx.Send(message, keyboard)
	}

	return ctx.Send(message)
}

func buildPaginationKeyboard(currentPage, totalPages int) *tg.ReplyMarkup {
	markup := &tg.ReplyMarkup{}
	var buttons []tg.Btn

	if currentPage > 1 {
		buttons = append(buttons, markup.Data("◀️ Previous", "list_page", fmt.Sprintf("%d", currentPage-1)))
	}

	buttons = append(buttons, markup.Data(fmt.Sprintf("%d / %d", currentPage, totalPages), "list_noop", ""))

	if currentPage < totalPages {
		buttons = append(buttons, markup.Data("Next ▶️", "list_page", fmt.Sprintf("%d", currentPage+1)))
	}

	markup.Inline(markup.Row(buttons...))
	return markup
}

func HandleListCallback(ctx tg.Context, packs *db.PackRepository, users *db.UserRepository) error {
	data := ctx.Callback().Data
	page := 1

	if parsedPage, err := fmt.Sscanf(data, "%d", &page); err != nil || parsedPage != 1 || page < 1 {
		page = 1
	}

	userID := ctx.Sender().ID
	lang := utils.GetUserLanguage(users, userID, ctx.Callback().Sender.LanguageCode)

	const pageSize = 5

	packList, total, err := packs.GetPaginated(userID, page, pageSize)
	if err != nil {
		utils.Logger("error", "Error getting packs for user", map[string]any{
			"userId": userID,
			"error":  err.Error(),
		})
		return ctx.Respond(&tg.CallbackResponse{Text: utils.T(lang, "error")})
	}

	if len(packList) == 0 {
		return ctx.Respond(&tg.CallbackResponse{Text: utils.T(lang, "list-page-empty")})
	}

	totalPages := (total + pageSize - 1) / pageSize
	message := utils.T(lang, "list-header-paginated", page, totalPages, total)

	startIndex := (page - 1) * pageSize
	for i, pack := range packList {
		orderNum := startIndex + i + 1
		message += utils.T(lang, "list-item", orderNum, pack.PackTitle, pack.PackType, pack.StickerCount, pack.PackLink)
	}

	if totalPages > 1 {
		keyboard := buildPaginationKeyboard(page, totalPages)
		ctx.Edit(message, keyboard)
	} else {
		ctx.Edit(message)
	}

	return ctx.Respond()
}

func HandleDeletePack(ctx tg.Context, relativeID int, packs *db.PackRepository, users *db.UserRepository) error {
	userID := ctx.Sender().ID
	lang := utils.GetUserLanguage(users, userID, ctx.Message().Sender.LanguageCode)

	pack, err := packs.GetByRelativeID(userID, relativeID)
	if err != nil {
		utils.Logger("error", "Error getting pack by relative ID", map[string]any{
			"relativeID": relativeID,
			"userId":     userID,
			"error":      err.Error(),
		})
		return ctx.Send(utils.T(lang, "error"))
	}

	if pack == nil {
		return ctx.Send(utils.T(lang, "delete-not-found"))
	}

	err = packs.Delete(pack.ID, userID)
	if err != nil {
		utils.Logger("error", "Error deleting pack", map[string]any{
			"packId": pack.ID,
			"userId": userID,
			"error":  err.Error(),
		})
		return ctx.Send(utils.T(lang, "delete-not-found"))
	}

	return ctx.Send(utils.T(lang, "delete-success"))
}
