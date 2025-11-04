package db

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		&Pack{},
		&User{},
		&SubscriptionPrice{},
		&UserSubscription{},
		&PaymentHistory{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	// Initialize default subscription prices
	repo := &Repository{db: db}
	if err := repo.initializeDefaultPrices(); err != nil {
		return nil, fmt.Errorf("failed to initialize default prices: %w", err)
	}

	return repo, nil
}

func (r *Repository) initializeDefaultPrices() error {
	defaultPrices := []SubscriptionPrice{
		{SubscriptionType: SubscriptionOneSteal, PriceStars: 10, Description: "Single emoji steal", Value: 1},
		{SubscriptionType: SubscriptionTenSteals, PriceStars: 80, Description: "10 emoji steals", Value: 10},
		{SubscriptionType: SubscriptionWeek, PriceStars: 50, Description: "1 week unlimited", Value: 7},
		{SubscriptionType: SubscriptionMonth, PriceStars: 150, Description: "1 month unlimited", Value: 30},
		{SubscriptionType: SubscriptionYear, PriceStars: 1200, Description: "1 year unlimited", Value: 365},
	}

	for _, price := range defaultPrices {
		// Use FirstOrCreate to avoid duplicate errors
		var existing SubscriptionPrice
		result := r.db.Where("subscription_type = ?", price.SubscriptionType).FirstOrCreate(&existing, price)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Pack methods

func (r *Repository) CreatePack(pack *Pack) error {
	result := r.db.Create(pack)
	return result.Error
}

func (r *Repository) GetPacksByUserID(userID int64) ([]Pack, error) {
	var packs []Pack
	result := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&packs)
	return packs, result.Error
}

func (r *Repository) GetPackByID(packID, userID int64) (*Pack, error) {
	var pack Pack
	result := r.db.Where("id = ? AND user_id = ?", packID, userID).First(&pack)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &pack, result.Error
}

func (r *Repository) DeletePack(packID, userID int64) error {
	result := r.db.Where("id = ? AND user_id = ?", packID, userID).Delete(&Pack{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("pack not found or not owned by user")
	}
	return nil
}

// User methods

func (r *Repository) UpsertUser(user *User) error {
	// Check if user exists
	var existing User
	result := r.db.Where("user_id = ?", user.UserID).First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new user
		return r.db.Create(user).Error
	}

	if result.Error != nil {
		return result.Error
	}

	// Update existing user
	return r.db.Model(&existing).Updates(map[string]interface{}{
		"username":      user.Username,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"language_code": user.LanguageCode,
		"last_seen_at":  time.Now(),
	}).Error
}

func (r *Repository) GetAllActiveUsers() ([]User, error) {
	var users []User
	result := r.db.Where("is_active = ?", true).Order("last_seen_at DESC").Find(&users)
	return users, result.Error
}

func (r *Repository) GetUserCount() (int, error) {
	var count int64
	result := r.db.Model(&User{}).Where("is_active = ?", true).Count(&count)
	return int(count), result.Error
}

func (r *Repository) GetPackCount() (int, error) {
	var count int64
	result := r.db.Model(&Pack{}).Count(&count)
	return int(count), result.Error
}

// Subscription methods

func (r *Repository) GetSubscriptionPrice(subType SubscriptionType) (*SubscriptionPrice, error) {
	var price SubscriptionPrice
	result := r.db.Where("subscription_type = ?", subType).First(&price)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &price, result.Error
}

func (r *Repository) GetAllSubscriptionPrices() ([]SubscriptionPrice, error) {
	var prices []SubscriptionPrice
	result := r.db.Order("price_stars ASC").Find(&prices)
	return prices, result.Error
}

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

func (r *Repository) CreateUserSubscription(sub *UserSubscription) error {
	result := r.db.Create(sub)
	return result.Error
}

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

func (r *Repository) DecrementSubscriptionCount(subID int64) error {
	result := r.db.Model(&UserSubscription{}).
		Where("id = ? AND remaining_count > 0", subID).
		UpdateColumn("remaining_count", gorm.Expr("remaining_count - 1"))
	return result.Error
}

func (r *Repository) CreatePaymentHistory(payment *PaymentHistory) error {
	result := r.db.Create(payment)
	return result.Error
}
