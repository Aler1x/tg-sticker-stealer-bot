package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Repository provides database access methods
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository instance and initializes the database
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

	// Initialize repository
	repo := &Repository{db: db}

	// Initialize default subscription prices
	if err := repo.initializeDefaultPrices(); err != nil {
		return nil, fmt.Errorf("failed to initialize default prices: %w", err)
	}

	return repo, nil
}

// initializeDefaultPrices creates default subscription prices if they don't exist
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

// Close closes the database connection
func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
