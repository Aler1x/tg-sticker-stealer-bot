package db

import (
	"time"

	"gorm.io/gorm"
)

type PackType string

const (
	PackTypeSticker PackType = "sticker"
	PackTypeEmoji   PackType = "emoji"
)

type Pack struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	UserID       int64     `gorm:"not null;index;uniqueIndex:idx_user_pack"`
	PackName     string    `gorm:"not null;uniqueIndex:idx_user_pack"`
	PackTitle    string    `gorm:"not null"`
	PackType     PackType  `gorm:"not null"`
	PackLink     string    `gorm:"not null"`
	StickerCount int       `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime;index"`
}

// TableName overrides the table name
func (Pack) TableName() string {
	return "packs"
}

type SubscriptionType string

const (
	SubscriptionOneSteal   SubscriptionType = "one_steal"
	SubscriptionTenSteals  SubscriptionType = "ten_steals"
	SubscriptionWeek       SubscriptionType = "week"
	SubscriptionMonth      SubscriptionType = "month"
	SubscriptionYear       SubscriptionType = "year"
)

type SubscriptionPrice struct {
	ID               int64            `gorm:"primaryKey;autoIncrement"`
	SubscriptionType SubscriptionType `gorm:"uniqueIndex;not null"`
	PriceStars       int              `gorm:"not null"`
	Description      string           `gorm:"not null"`
	Value            int              `gorm:"not null"`
	UpdatedAt        time.Time        `gorm:"autoUpdateTime"`
}

func (SubscriptionPrice) TableName() string {
	return "subscription_prices"
}

type UserSubscription struct {
	ID               int64            `gorm:"primaryKey;autoIncrement"`
	UserID           int64            `gorm:"not null;index"`
	SubscriptionType SubscriptionType `gorm:"not null"`
	RemainingCount   *int             `gorm:"default:null"`
	ExpiresAt        *time.Time       `gorm:"index;default:null"`
	CreatedAt        time.Time        `gorm:"autoCreateTime"`
}

func (UserSubscription) TableName() string {
	return "user_subscriptions"
}

type PaymentHistory struct {
	ID              int64            `gorm:"primaryKey;autoIncrement"`
	UserID          int64            `gorm:"not null;index"`
	SubscriptionType SubscriptionType `gorm:"not null"`
	PriceStars      int              `gorm:"not null"`
	PaymentChargeID *string          `gorm:"default:null"`
	PaymentProvider *string          `gorm:"default:null"`
	Status          string           `gorm:"not null;index"`
	CreatedAt       time.Time        `gorm:"autoCreateTime"`
}

func (PaymentHistory) TableName() string {
	return "payment_history"
}

