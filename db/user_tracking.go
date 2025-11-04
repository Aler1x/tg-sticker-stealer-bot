package db

import "time"

type User struct {
	UserID       int64     `gorm:"primaryKey;column:user_id"`
	Username     string    `gorm:"column:username"`
	FirstName    string    `gorm:"column:first_name"`
	LastName     string    `gorm:"column:last_name"`
	LanguageCode string    `gorm:"column:language_code"`
	IsActive     bool      `gorm:"column:is_active;default:1;index"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	LastSeenAt   time.Time `gorm:"column:last_seen_at;autoUpdateTime;index"`
}

func (User) TableName() string {
	return "users"
}

