package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"tg-sticker-stiller-bot/types"
	"tg-sticker-stiller-bot/utils"

	"github.com/google/uuid"
	tg "gopkg.in/telebot.v4"
)

const TempDir = "./data/temp"

func DownloadFile(bot *tg.Bot, file tg.File) (string, error) {
	reader, err := bot.File(&file)
	if err != nil {
		return "", fmt.Errorf("failed to get file: %w", err)
	}
	defer reader.Close()

	var extension string
	if sticker, ok := interface{}(&file).(*tg.Sticker); ok {
		extension = getFileExtension(*sticker)
	} else {
		extension = "webp"
	}

	filename := fmt.Sprintf("%s.%s", uuid.New().String(), extension)
	filePath := filepath.Join(TempDir, filename)

	if err := utils.EnsureTempDir(); err != nil {
		return "", err
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, reader); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

func DownloadSticker(bot *tg.Bot, sticker tg.Sticker) (string, error) {
	reader, err := bot.File(&sticker.File)
	if err != nil {
		return "", fmt.Errorf("failed to get file: %w", err)
	}
	defer reader.Close()

	extension := getFileExtension(sticker)
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), extension)
	filePath := filepath.Join(TempDir, filename)

	if err := utils.EnsureTempDir(); err != nil {
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

func DownloadAllStickers(bot *tg.Bot, stickers []tg.Sticker) []types.DownloadedSticker {
	var wg sync.WaitGroup
	results := make([]*types.DownloadedSticker, len(stickers))

	for i, sticker := range stickers {
		wg.Add(1)
		go func(index int, s tg.Sticker) {
			defer wg.Done()

			filePath, err := DownloadSticker(bot, s)
			if err != nil {
				utils.Logger("warn", "Failed to download sticker, skipping", map[string]any{
					"fileId": s.FileID,
					"error":  err.Error(),
				})
				results[index] = nil
				return
			}

			results[index] = &types.DownloadedSticker{
				Path:    filePath,
				Sticker: s,
			}
		}(i, sticker)
	}

	wg.Wait()

	downloaded := []types.DownloadedSticker{}
	for _, result := range results {
		if result != nil {
			downloaded = append(downloaded, *result)
		}
	}

	return downloaded
}

func getFileExtension(item tg.Sticker) string {
	if item.Animated {
		return "tgs"
	}
	if item.Video {
		return "webm"
	}
	return "webp"
}
