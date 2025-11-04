package utils

import (
	"fmt"
	"os"
	"sync"
	"time"

	tg "gopkg.in/telebot.v4"
)

const (
	MaxRetries = 3
	RetryDelay = 2 * time.Second
	TempDir    = "./data/temp"
)

func EnsureTempDir() error {
	if err := os.MkdirAll(TempDir, 0755); err != nil {
		Logger("error", "Failed to create temp directory", map[string]any{"error": err.Error()})
		return err
	}
	return nil
}

func WithRetry[T any](fn func() (T, error)) (T, error) {
	var result T
	var err error

	for attempt := range MaxRetries {
		result, err = fn()
		if err == nil {
			return result, nil
		}

		if _, ok := err.(*BotError); ok {
			return result, err
		}

		if attempt < MaxRetries-1 {
			Logger("warn", "Retry attempt failed", map[string]any{
				"attempt": attempt + 1,
				"error":   err.Error(),
				"delay":   RetryDelay.String(),
			})
			time.Sleep(RetryDelay)
		}
	}

	return result, fmt.Errorf("max retries exceeded: %w", err)
}

func GenerateSetName(name, botname string) string {
	return fmt.Sprintf("%s_by_%s", name, botname)
}

func GetStickerFormat(sticker tg.Sticker) string {
	if sticker.Animated {
		return tg.StickerAnimated
	}
	if sticker.Video {
		return tg.StickerVideo
	}
	return tg.StickerStatic
}

func CleanupFiles(filePaths []string) {
	var wg sync.WaitGroup

	for _, filePath := range filePaths {
		wg.Add(1)
		go func(fp string) {
			defer wg.Done()

			if err := os.Remove(fp); err != nil {
				Logger("warn", "Failed to delete file", map[string]any{
					"filePath": fp,
					"error":    err.Error(),
				})
			}
		}(filePath)
	}

	wg.Wait()
}
