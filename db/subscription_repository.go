package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// GetSubscriptionPrice retrieves the pricing for a specific subscription type
func (r *Repository) GetSubscriptionPrice(subType SubscriptionType) (*SubscriptionPrice, error) {
	var price SubscriptionPrice
	result := r.db.Where("subscription_type = ?", subType).First(&price)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &price, result.Error
}

// GetAllSubscriptionPrices retrieves all subscription prices ordered by price
func (r *Repository) GetAllSubscriptionPrices() ([]SubscriptionPrice, error) {
	var prices []SubscriptionPrice
	result := r.db.Order("price_stars ASC").Find(&prices)
	return prices, result.Error
}

// UpdateSubscriptionPrice updates the price for a specific subscription type
func (r *Repository) UpdateSubscriptionPrice(subType SubscriptionType, priceStars int) error {
	result := r.db.Model(&SubscriptionPrice{}).
		Where("subscription_type = ?", subType).
		Updates(map[string]interface{}{
			"price_stars": priceStars,
			"updated_at":  time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription type not found")
	}

	return nil
}

// CreateUserSubscription creates a new subscription for a user
func (r *Repository) CreateUserSubscription(sub *UserSubscription) error {
	result := r.db.Create(sub)
	return result.Error
}

// GetActiveSubscription retrieves the active subscription for a user
func (r *Repository) GetActiveSubscription(userID int64) (*UserSubscription, error) {
	var sub UserSubscription
	result := r.db.Where("user_id = ?", userID).
		Where("(remaining_count IS NOT NULL AND remaining_count > 0) OR (expires_at IS NOT NULL AND expires_at > ?)", time.Now()).
		Order("created_at DESC").
		First(&sub)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &sub, result.Error
}

// DecrementSubscriptionCount decrements the remaining count for a subscription
func (r *Repository) DecrementSubscriptionCount(subID int64) error {
	result := r.db.Model(&UserSubscription{}).
		Where("id = ? AND remaining_count > 0", subID).
		UpdateColumn("remaining_count", gorm.Expr("remaining_count - 1"))
	return result.Error
}
