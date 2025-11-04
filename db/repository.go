package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err := db.Exec(Schema); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) CreatePack(pack *Pack) error {
	query := `
		INSERT INTO packs (user_id, pack_name, pack_title, pack_type, pack_link, sticker_count)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, pack.UserID, pack.PackName, pack.PackTitle, pack.PackType, pack.PackLink, pack.StickerCount)
	if err != nil {
		return fmt.Errorf("failed to create pack: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	pack.ID = id
	return nil
}

func (r *Repository) GetPacksByUserID(userID int64) ([]Pack, error) {
	query := `
		SELECT id, user_id, pack_name, pack_title, pack_type, pack_link, sticker_count, created_at
		FROM packs
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query packs: %w", err)
	}
	defer rows.Close()

	var packs []Pack
	for rows.Next() {
		var pack Pack
		err := rows.Scan(&pack.ID, &pack.UserID, &pack.PackName, &pack.PackTitle, &pack.PackType, &pack.PackLink, &pack.StickerCount, &pack.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pack: %w", err)
		}
		packs = append(packs, pack)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating packs: %w", err)
	}

	return packs, nil
}

func (r *Repository) GetPackByID(packID, userID int64) (*Pack, error) {
	query := `
		SELECT id, user_id, pack_name, pack_title, pack_type, pack_link, sticker_count, created_at
		FROM packs
		WHERE id = ? AND user_id = ?
	`
	var pack Pack
	err := r.db.QueryRow(query, packID, userID).Scan(
		&pack.ID, &pack.UserID, &pack.PackName, &pack.PackTitle, &pack.PackType, &pack.PackLink, &pack.StickerCount, &pack.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pack: %w", err)
	}

	return &pack, nil
}

func (r *Repository) DeletePack(packID, userID int64) error {
	query := `DELETE FROM packs WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, packID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete pack: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pack not found or not owned by user")
	}

	return nil
}

func (r *Repository) UpsertUser(user *User) error {
	query := `
		INSERT INTO users (user_id, username, first_name, last_name, language_code, last_seen_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id) DO UPDATE SET
			username = excluded.username,
			first_name = excluded.first_name,
			last_name = excluded.last_name,
			language_code = excluded.language_code,
			last_seen_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, user.UserID, user.Username, user.FirstName, user.LastName, user.LanguageCode)
	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}
	return nil
}

func (r *Repository) GetAllActiveUsers() ([]User, error) {
	query := `
		SELECT user_id, username, first_name, last_name, language_code, is_active, created_at, last_seen_at
		FROM users
		WHERE is_active = 1
		ORDER BY last_seen_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.LanguageCode, &user.IsActive, &user.CreatedAt, &user.LastSeenAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *Repository) GetUserCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE is_active = 1`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

func (r *Repository) GetPackCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM packs`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count packs: %w", err)
	}
	return count, nil
}

// Subscription methods

func (r *Repository) GetSubscriptionPrice(subType SubscriptionType) (*SubscriptionPrice, error) {
	query := `
		SELECT id, subscription_type, price_stars, description, value, updated_at
		FROM subscription_prices
		WHERE subscription_type = ?
	`
	var price SubscriptionPrice
	err := r.db.QueryRow(query, subType).Scan(
		&price.ID, &price.SubscriptionType, &price.PriceStars, &price.Description, &price.Value, &price.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription price: %w", err)
	}
	return &price, nil
}

func (r *Repository) GetAllSubscriptionPrices() ([]SubscriptionPrice, error) {
	query := `
		SELECT id, subscription_type, price_stars, description, value, updated_at
		FROM subscription_prices
		ORDER BY price_stars ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscription prices: %w", err)
	}
	defer rows.Close()

	var prices []SubscriptionPrice
	for rows.Next() {
		var price SubscriptionPrice
		err := rows.Scan(&price.ID, &price.SubscriptionType, &price.PriceStars, &price.Description, &price.Value, &price.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription price: %w", err)
		}
		prices = append(prices, price)
	}

	return prices, rows.Err()
}

func (r *Repository) UpdateSubscriptionPrice(subType SubscriptionType, priceStars int) error {
	query := `
		UPDATE subscription_prices
		SET price_stars = ?, updated_at = CURRENT_TIMESTAMP
		WHERE subscription_type = ?
	`
	result, err := r.db.Exec(query, priceStars, subType)
	if err != nil {
		return fmt.Errorf("failed to update subscription price: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription type not found")
	}

	return nil
}

func (r *Repository) CreateUserSubscription(sub *UserSubscription) error {
	query := `
		INSERT INTO user_subscriptions (user_id, subscription_type, remaining_count, expires_at)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, sub.UserID, sub.SubscriptionType, sub.RemainingCount, sub.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	sub.ID = id
	return nil
}

func (r *Repository) GetActiveSubscription(userID int64) (*UserSubscription, error) {
	query := `
		SELECT id, user_id, subscription_type, remaining_count, expires_at, created_at
		FROM user_subscriptions
		WHERE user_id = ?
		AND (
			(remaining_count IS NOT NULL AND remaining_count > 0)
			OR (expires_at IS NOT NULL AND expires_at > CURRENT_TIMESTAMP)
		)
		ORDER BY created_at DESC
		LIMIT 1
	`
	var sub UserSubscription
	err := r.db.QueryRow(query, userID).Scan(
		&sub.ID, &sub.UserID, &sub.SubscriptionType, &sub.RemainingCount, &sub.ExpiresAt, &sub.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscription: %w", err)
	}
	return &sub, nil
}

func (r *Repository) DecrementSubscriptionCount(subID int64) error {
	query := `
		UPDATE user_subscriptions
		SET remaining_count = remaining_count - 1
		WHERE id = ? AND remaining_count > 0
	`
	_, err := r.db.Exec(query, subID)
	if err != nil {
		return fmt.Errorf("failed to decrement subscription count: %w", err)
	}
	return nil
}

func (r *Repository) CreatePaymentHistory(payment *PaymentHistory) error {
	query := `
		INSERT INTO payment_history (user_id, subscription_type, price_stars, payment_charge_id, payment_provider, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, payment.UserID, payment.SubscriptionType, payment.PriceStars, payment.PaymentChargeID, payment.PaymentProvider, payment.Status)
	if err != nil {
		return fmt.Errorf("failed to create payment history: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	payment.ID = id
	return nil
}
