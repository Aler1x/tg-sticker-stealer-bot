package services

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"tg-sticker-stiller-bot/types"
	"tg-sticker-stiller-bot/utils"

	"github.com/google/uuid"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
	tg "gopkg.in/telebot.v4"
)

const (
	maxStickerDimension = 512
	emojiDimension      = 100
)

func PrepareStickerForSet(filePath string, sticker tg.Sticker, targetType types.StickerType) (string, error) {
	if sticker.Animated || sticker.Video {
		return filePath, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open sticker file: %w", err)
	}
	defer file.Close()

	srcImg, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode sticker: %w", err)
	}

	bounds := srcImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if !stickerNeedsResize(width, height, targetType) {
		return filePath, nil
	}

	var dstImg *image.RGBA
	if targetType == types.StickerTypeEmoji {
		dstImg = resizeToSquareCanvas(srcImg, emojiDimension)
	} else {
		dstImg = resizeToStickerDimensions(srcImg)
	}

	if err := utils.EnsureTempDir(); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.png", uuid.New().String())
	outputPath := filepath.Join(utils.TempDir, filename)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create prepared sticker file: %w", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, dstImg); err != nil {
		os.Remove(outputPath)
		return "", fmt.Errorf("failed to encode prepared sticker: %w", err)
	}

	return outputPath, nil
}

func stickerNeedsResize(width, height int, targetType types.StickerType) bool {
	if targetType == types.StickerTypeEmoji {
		return width != emojiDimension || height != emojiDimension
	}

	return !isValidStickerDimensions(width, height)
}

func isValidStickerDimensions(width, height int) bool {
	return (width == maxStickerDimension && height <= maxStickerDimension) ||
		(height == maxStickerDimension && width <= maxStickerDimension)
}

func resizeToStickerDimensions(src image.Image) *image.RGBA {
	bounds := src.Bounds()
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
	draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), src, bounds, draw.Over, nil)
	return dstImg
}

func resizeToSquareCanvas(src image.Image, size int) *image.RGBA {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	scale := min(float64(size)/float64(width), float64(size)/float64(height))
	newWidth := int(float64(width) * scale)
	newHeight := int(float64(height) * scale)

	scaled := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(scaled, scaled.Bounds(), src, bounds, draw.Over, nil)

	canvas := image.NewRGBA(image.Rect(0, 0, size, size))
	offsetX := (size - newWidth) / 2
	offsetY := (size - newHeight) / 2
	draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+newWidth, offsetY+newHeight), scaled, image.Point{}, draw.Over)

	return canvas
}
