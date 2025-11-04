package db

import "time"

type PackType string

const (
	PackTypeSticker PackType = "sticker"
	PackTypeEmoji   PackType = "emoji"
)

type Pack struct {
	ID           int64     `db:"id"`
	UserID       int64     `db:"user_id"`
	PackName     string    `db:"pack_name"`
	PackTitle    string    `db:"pack_title"`
	PackType     PackType  `db:"pack_type"`
	PackLink     string    `db:"pack_link"`
	StickerCount int       `db:"sticker_count"`
	CreatedAt    time.Time `db:"created_at"`
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
	ID               int64            `db:"id"`
	SubscriptionType SubscriptionType `db:"subscription_type"`
	PriceStars       int              `db:"price_stars"`
	Description      string           `db:"description"`
	Value            int              `db:"value"`
	UpdatedAt        time.Time        `db:"updated_at"`
}

type UserSubscription struct {
	ID               int64            `db:"id"`
	UserID           int64            `db:"user_id"`
	SubscriptionType SubscriptionType `db:"subscription_type"`
	RemainingCount   *int             `db:"remaining_count"`
	ExpiresAt        *time.Time       `db:"expires_at"`
	CreatedAt        time.Time        `db:"created_at"`
}

type PaymentHistory struct {
	ID                int64            `db:"id"`
	UserID            int64            `db:"user_id"`
	SubscriptionType  SubscriptionType `db:"subscription_type"`
	PriceStars        int              `db:"price_stars"`
	PaymentChargeID   *string          `db:"payment_charge_id"`
	PaymentProvider   *string          `db:"payment_provider"`
	Status            string           `db:"status"`
	CreatedAt         time.Time        `db:"created_at"`
}

