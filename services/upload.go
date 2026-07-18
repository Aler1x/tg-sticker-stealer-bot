package services

import (
	"fmt"
	"strings"
	"tg-sticker-stiller-bot/types"
	"tg-sticker-stiller-bot/utils"
	"time"

	tg "gopkg.in/telebot.v4"
)

type ProgressCallback func(current, total int)

func CreateStickerSet(bot *tg.Bot, userID int64, botname string, title string, stickers []tg.Sticker, stickerType types.StickerType, progressCallback ProgressCallback) (string, error) {
	downloadedStickers := DownloadAllStickers(bot, stickers)
	if len(downloadedStickers) == 0 {
		utils.Logger("error", "No stickers could be downloaded", map[string]any{"userId": userID})
		return "", fmt.Errorf("no stickers could be downloaded")
	}

	filePaths := make([]string, 0, len(downloadedStickers)*2)
	preparedStickers := make([]types.DownloadedSticker, 0, len(downloadedStickers))

	for _, ds := range downloadedStickers {
		filePaths = append(filePaths, ds.Path)

		preparedPath, err := PrepareStickerForSet(ds.Path, ds.Sticker, stickerType)
		if err != nil {
			utils.Logger("warn", "Failed to prepare sticker for set, skipping", map[string]any{
				"fileId": ds.Sticker.FileID,
				"error":  err.Error(),
			})
			continue
		}

		if preparedPath != ds.Path {
			filePaths = append(filePaths, preparedPath)
		}

		preparedStickers = append(preparedStickers, types.DownloadedSticker{
			Path:    preparedPath,
			Sticker: ds.Sticker,
		})
	}
	defer utils.CleanupFiles(filePaths)

	if len(preparedStickers) == 0 {
		utils.Logger("error", "No stickers could be prepared", map[string]any{"userId": userID})
		return "", fmt.Errorf("no stickers could be prepared")
	}

	downloadedStickers = preparedStickers

	normalizedName := utils.NormalizePackName(title)
	setName := utils.GenerateSetName(normalizedName, botname)

	user := &tg.User{ID: userID}

	var telegramStickerType tg.StickerSetType
	var packLink string

	if stickerType == types.StickerTypeEmoji {
		telegramStickerType = tg.StickerCustomEmoji
		packLink = fmt.Sprintf("https://t.me/addemoji/%s", setName)
	} else {
		telegramStickerType = tg.StickerRegular
		packLink = fmt.Sprintf("https://t.me/addstickers/%s", setName)
	}

	firstSticker := downloadedStickers[0]
	emoji := firstSticker.Sticker.Emoji
	if emoji == "" {
		emoji = "😀"
	}

	err := utils.CallWithFloodRetry(func() error {
		stickerSet := &tg.StickerSet{
			Type:  telegramStickerType,
			Name:  setName,
			Title: title,
			Input: []tg.InputSticker{{
				File:     tg.FromDisk(firstSticker.Path),
				Format:   utils.GetStickerFormat(firstSticker.Sticker),
				Emojis:   []string{emoji},
				Keywords: []string{},
			}},
		}
		return bot.CreateStickerSet(user, stickerSet)
	})
	if err != nil {
		if isNameTakenError(err) {
			utils.Logger("warn", "Sticker set name already exists", map[string]any{
				"title":  title,
				"userId": userID,
			})
			return "", utils.NewBotError(
				fmt.Sprintf("Sticker set name already exists: %s", title),
				"name-taken",
				"STICKER_SET_NAME_TAKEN",
			)
		}
		utils.Logger("error", "Failed to create sticker set", map[string]any{
			"userId": userID,
			"error":  err.Error(),
		})
		return "", err
	}

	totalStickers := len(downloadedStickers)
	addedCount := 1

	if progressCallback != nil && totalStickers > 1 {
		progressCallback(1, totalStickers)
	}

	for i := 1; i < totalStickers; i++ {
		stickerData := downloadedStickers[i]
		emoji := stickerData.Sticker.Emoji
		if emoji == "" {
			emoji = "😀"
		}

		err := utils.CallWithFloodRetry(func() error {
			return bot.AddStickerToSet(user, setName, tg.InputSticker{
				File:     tg.FromDisk(stickerData.Path),
				Format:   utils.GetStickerFormat(stickerData.Sticker),
				Emojis:   []string{emoji},
				Keywords: []string{},
			})
		})
		if err != nil {
			utils.Logger("warn", "Failed to add sticker to set", map[string]any{
				"current": i + 1,
				"total":   totalStickers,
				"error":   err.Error(),
			})
		} else {
			addedCount++
		}

		if progressCallback != nil {
			if (i+1)%10 == 0 || i+1 == totalStickers {
				progressCallback(i+1, totalStickers)
			}
		}

		if i < totalStickers-1 {
			time.Sleep(utils.StickerUploadPace())
		}
	}

	if addedCount == 0 {
		return "", fmt.Errorf("no stickers could be added to the set")
	}

	if addedCount < totalStickers {
		utils.Logger("warn", "Sticker set created with missing items", map[string]any{
			"userId": userID,
			"added":  addedCount,
			"total":  totalStickers,
		})
	}

	return packLink, nil
}

func isNameTakenError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "name is already occupied") || strings.Contains(errStr, "409")
}
