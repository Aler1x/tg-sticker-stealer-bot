package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) Upsert(user *User) error {
	now := time.Now()
	result := r.db.Where("user_id = ?", user.UserID).First(&User{})

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user.CreatedAt = now
		user.LastSeenAt = now
		user.IsActive = true
		if err := r.db.Create(user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	}

	if result.Error != nil {
		return fmt.Errorf("failed to check user: %w", result.Error)
	}

	if err := r.db.Model(&User{}).Where("user_id = ?", user.UserID).Updates(map[string]any{
		"username":      user.Username,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"language_code": user.LanguageCode,
		"last_seen_at":  now,
	}).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetAllActive() ([]User, error) {
	var users []User
	if err := r.db.Where("is_active = ?", true).Order("last_seen_at DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	return users, nil
}

func (r *UserRepository) Count() (int, error) {
	var count int64
	if err := r.db.Model(&User{}).Where("is_active = ?", true).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return int(count), nil
}

func (r *UserRepository) GetByID(userID int64) (*User, error) {
	var user User
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) SetDefaultAction(userID int64, action DefaultAction) error {
	if err := r.db.Model(&User{}).Where("user_id = ?", userID).Update("default_action", action).Error; err != nil {
		return fmt.Errorf("failed to update default action: %w", err)
	}
	return nil
}
