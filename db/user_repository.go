package db

import (
	"time"

	"gorm.io/gorm"
)

// UpsertUser creates a new user or updates an existing one
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

// GetAllActiveUsers retrieves all active users
func (r *Repository) GetAllActiveUsers() ([]User, error) {
	var users []User
	result := r.db.Where("is_active = ?", true).Order("last_seen_at DESC").Find(&users)
	return users, result.Error
}

// GetUserCount returns the total number of active users
func (r *Repository) GetUserCount() (int, error) {
	var count int64
	result := r.db.Model(&User{}).Where("is_active = ?", true).Count(&count)
	return int(count), result.Error
}
