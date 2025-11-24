package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"tg-sticker-stiller-bot/utils"

	"github.com/google/uuid"
	tg "gopkg.in/telebot.v4"
)

func CreateStickerZip(bot *tg.Bot, stickers []tg.Sticker, packName string, progressCallback ProgressCallback) (string, error) {
	downloadedStickers := DownloadAllStickers(bot, stickers)
	if len(downloadedStickers) == 0 {
		return "", fmt.Errorf("no stickers could be downloaded")
	}

	filePaths := make([]string, len(downloadedStickers))
	for i, ds := range downloadedStickers {
		filePaths[i] = ds.Path
	}
	defer utils.CleanupFiles(filePaths)

	zipFileName := fmt.Sprintf("%s_%s.zip", packName, uuid.New().String()[:8])
	zipFilePath := filepath.Join(TempDir, zipFileName)

	if err := utils.EnsureTempDir(); err != nil {
		return "", err
	}

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	total := len(downloadedStickers)
	for i, ds := range downloadedStickers {
		ext := getFileExtension(ds.Sticker)
		stickerFileName := fmt.Sprintf("%03d.%s", i+1, ext)

		f, err := os.Open(ds.Path)
		if err != nil {
			utils.Logger("warn", "Failed to open sticker file", map[string]any{"path": ds.Path, "error": err.Error()})
			continue
		}

		w, err := zipWriter.Create(stickerFileName)
		if err != nil {
			f.Close()
			utils.Logger("warn", "Failed to create zip entry", map[string]any{"name": stickerFileName, "error": err.Error()})
			continue
		}

		if _, err := io.Copy(w, f); err != nil {
			f.Close()
			utils.Logger("warn", "Failed to write to zip", map[string]any{"name": stickerFileName, "error": err.Error()})
			continue
		}
		f.Close()

		if progressCallback != nil && ((i+1)%10 == 0 || i+1 == total) {
			progressCallback(i+1, total)
		}
	}

	return zipFilePath, nil
}
