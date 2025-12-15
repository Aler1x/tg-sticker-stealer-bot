package services

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"tg-sticker-stiller-bot/utils"

	"github.com/google/uuid"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
	tg "gopkg.in/telebot.v4"
)

const (
	maxStickerDimension = 512
)

func ConvertImageToSticker(bot *tg.Bot, photo *tg.Photo) (string, error) {
	if photo == nil || photo.FileID == "" {
		return "", fmt.Errorf("no photo provided")
	}

	reader, err := bot.File(&photo.File)
	if err != nil {
		return "", fmt.Errorf("failed to get photo file: %w", err)
	}
	defer reader.Close()

	srcImg, _, err := image.Decode(reader)
	if err != nil {
		return "", utils.NewBotError("failed to decode image", "unsupported-format", "UNSUPPORTED_FORMAT")
	}

	bounds := srcImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var newWidth, newHeight int
	if width > height {
		newWidth = maxStickerDimension
		newHeight = int(float64(height) * float64(maxStickerDimension) / float64(width))
	} else {
		newHeight = maxStickerDimension
		newWidth = int(float64(width) * float64(maxStickerDimension) / float64(height))
	}

	dstImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, bounds, draw.Over, nil)

	if err := utils.EnsureTempDir(); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.png", uuid.New().String())
	filePath := filepath.Join(utils.TempDir, filename)

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, dstImg); err != nil {
		return "", fmt.Errorf("failed to encode png: %w", err)
	}

	return filePath, nil
}

func ConvertStickerToImage(bot *tg.Bot, sticker *tg.Sticker) (string, error) {
	if sticker == nil || sticker.FileID == "" {
		return "", fmt.Errorf("no sticker provided")
	}

	if sticker.Animated || sticker.Video {
		return "", fmt.Errorf("animated/video stickers cannot be converted to image")
	}

	reader, err := bot.File(&sticker.File)
	if err != nil {
		return "", fmt.Errorf("failed to get sticker file: %w", err)
	}
	defer reader.Close()

	srcImg, _, err := image.Decode(reader)
	if err != nil {
		return "", fmt.Errorf("failed to decode sticker: %w", err)
	}

	if err := utils.EnsureTempDir(); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.png", uuid.New().String())
	filePath := filepath.Join(utils.TempDir, filename)

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, srcImg); err != nil {
		return "", fmt.Errorf("failed to encode png: %w", err)
	}

	return filePath, nil
}
