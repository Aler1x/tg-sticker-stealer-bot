package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	PackTypeStickers = "stickers"
	PackTypeEmojis   = "emojis"
)

type PackCreation struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	UserID      int64     `gorm:"column:user_id;not null;index"`
	PackLink    string    `gorm:"column:pack_link;not null"`
	PackType    string    `gorm:"column:pack_type;not null;index"`
	ItemsAmount int       `gorm:"column:items_amount;not null;default:0"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime;index"`
}

type PackCreationRepository struct {
	db *gorm.DB
}

func (r *PackCreationRepository) Record(userID int64, packLink, packType string, itemsAmount int) error {
	row := &PackCreation{
		ID:          uuid.New(),
		UserID:      userID,
		PackLink:    packLink,
		PackType:    packType,
		ItemsAmount: itemsAmount,
	}
	if err := r.db.Create(row).Error; err != nil {
		return fmt.Errorf("failed to record pack creation: %w", err)
	}
	return nil
}

func (r *PackCreationRepository) CountStickerPacks() (int64, error) {
	var n int64
	err := r.db.Model(&PackCreation{}).
		Where("pack_type IN ?", []string{PackTypeStickers, "regular"}).
		Count(&n).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count sticker pack creations: %w", err)
	}
	return n, nil
}

func (r *PackCreationRepository) CountEmojiPacks() (int64, error) {
	var n int64
	err := r.db.Model(&PackCreation{}).
		Where("pack_type IN ?", []string{PackTypeEmojis, "custom_emoji"}).
		Count(&n).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count emoji pack creations: %w", err)
	}
	return n, nil
}
