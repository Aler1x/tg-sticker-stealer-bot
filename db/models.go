package db

import "time"

type PackType string

const (
	PackTypeSticker PackType = "sticker"
	PackTypeEmoji   PackType = "emoji"
)

type Pack struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64     `gorm:"column:user_id;not null;uniqueIndex:idx_user_pack"`
	PackName     string    `gorm:"column:pack_name;not null;uniqueIndex:idx_user_pack"`
	PackTitle    string    `gorm:"column:pack_title;not null"`
	PackType     PackType  `gorm:"column:pack_type;not null"`
	PackLink     string    `gorm:"column:pack_link;not null"`
	StickerCount int       `gorm:"column:sticker_count;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;index"`
}

type User struct {
	UserID       int64     `gorm:"column:user_id;primaryKey"`
	Username     string    `gorm:"column:username"`
	FirstName    string    `gorm:"column:first_name"`
	LastName     string    `gorm:"column:last_name"`
	LanguageCode string    `gorm:"column:language_code"`
	IsActive     bool      `gorm:"column:is_active;default:true;index"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	LastSeenAt   time.Time `gorm:"column:last_seen_at;autoCreateTime;index"`
}
