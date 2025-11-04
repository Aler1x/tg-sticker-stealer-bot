package handlers

import (
	"fmt"
	"sort"
	"tg-sticker-stiller-bot/db"
	"tg-sticker-stiller-bot/services"
	"tg-sticker-stiller-bot/types"
	"tg-sticker-stiller-bot/utils"

	tg "gopkg.in/telebot.v4"
)

func HandlePack(ctx tg.Context, packName string, packType types.StickerType, bot *tg.Bot, sessions *services.SessionStore, repo *db.Repository) error {
	lang := ctx.Message().Sender.LanguageCode
	userID := ctx.Sender().ID

	// Check subscription for emoji packs only
	if packType == types.StickerTypeEmoji {
		sub, err := CheckSubscription(userID, repo)
		if err != nil {
			utils.LogError("Failed to check subscription", err)
			return ctx.Send("Error checking subscription status. Please try again.")
		}

		if sub == nil {
			return ctx.Send(
				"🔒 *Emoji stealing requires a subscription!*\n\n" +
					"Get access to steal custom emoji packs by choosing a subscription plan.\n\n" +
					"Use /subscribe to view available options!",
				tg.ModeMarkdown,
			)
		}
	}

	var stickerSet *types.StickerSet
	var err error

	if packType == types.StickerTypeEmoji {
		emojiSet, fetchErr := services.FetchEmojiSet(bot, packName)
		if fetchErr != nil {
			utils.Logger("error", "Error fetching emoji pack", map[string]any{
				"packName": packName,
				"error":    fetchErr.Error(),
			})
			return ctx.Send(utils.T(lang, "error"))
		}
		stickerSet = &types.StickerSet{
			Name:     emojiSet.Name,
			Title:    emojiSet.Title,
			Stickers: emojiSet.Stickers,
		}
	} else {
		stickerSet, err = services.FetchStickerSet(bot, packName)
		if err != nil {
			utils.Logger("error", "Error fetching sticker pack", map[string]any{
				"packName": packName,
				"error":    err.Error(),
			})
			return ctx.Send(utils.T(lang, "error"))
		}
	}

	packTypeKey := "pack-type"
	if packType == types.StickerTypeEmoji {
		packTypeKey = "emoji-type"
	}

	ctx.Send(utils.T(lang, "pack-stats", utils.T(lang, packTypeKey), stickerSet.Title, len(stickerSet.Stickers)))

	sessions.Set(userID, &services.Session{
		State:         services.StateWaitingForPackName,
		Title:         stickerSet.Title,
		OriginalItems: stickerSet.Stickers,
		Name:          packName,
		PackType:      packType,
	})

	return nil
}

func HandlePackNameInput(ctx tg.Context, userInput string, bot *tg.Bot, sessions *services.SessionStore, repo *db.Repository) error {
	lang := ctx.Message().Sender.LanguageCode
	userID := ctx.Sender().ID

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

	packLink, err := services.CreateStickerSet(bot, userID, bot.Me.Username, userInput, session.OriginalItems, session.PackType, repo, progressCallback)
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

	// Consume subscription for emoji packs
	if session.PackType == types.StickerTypeEmoji {
		sub, err := CheckSubscription(userID, repo)
		if err != nil {
			utils.LogError("Failed to check subscription after creation", err)
		} else if sub != nil {
			if err := ConsumeSubscription(sub, repo); err != nil {
				utils.LogError("Failed to consume subscription", err)
			} else {
				utils.Logger("info", "Subscription consumed", map[string]any{
					"userId":           userID,
					"subscriptionType": sub.SubscriptionType,
				})
			}
		}
	}

	if progressMsg != nil {
		ctx.Bot().Delete(progressMsg)
	}

	ctx.Send(utils.T(lang, "success", utils.T(lang, packTypeKey), packLink))
	sessions.Clear(userID)
	return nil
}

func HandleListPacks(ctx tg.Context, repo *db.Repository) error {
	lang := ctx.Message().Sender.LanguageCode
	userID := ctx.Sender().ID

	packs, err := repo.GetPacksByUserID(userID)
	if err != nil {
		utils.Logger("error", "Error getting packs for user", map[string]any{
			"userId": userID,
			"error":  err.Error(),
		})
		return ctx.Send(utils.T(lang, "error"))
	}

	if len(packs) == 0 {
		return ctx.Send(utils.T(lang, "list-empty"))
	}

	sort.Slice(packs, func(i, j int) bool {
		return packs[i].ID < packs[j].ID
	})

	message := utils.T(lang, "list-header")
	for _, pack := range packs {
		message += utils.T(lang, "list-item", pack.ID, pack.PackTitle, pack.PackType, pack.StickerCount, pack.PackLink)
	}

	return ctx.Send(message)
}

func HandleDeletePack(ctx tg.Context, packID int64, repo *db.Repository) error {
	lang := ctx.Message().Sender.LanguageCode
	userID := ctx.Sender().ID

	err := repo.DeletePack(packID, userID)
	if err != nil {
		utils.Logger("error", "Error deleting pack", map[string]any{
			"packId": packID,
			"userId": userID,
			"error":  err.Error(),
		})
		return ctx.Send(utils.T(lang, "delete-not-found"))
	}

	return ctx.Send(utils.T(lang, "delete-success"))
}
