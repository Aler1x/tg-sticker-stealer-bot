package db

import "time"

type DefaultAction string

const (
	DefaultActionCopy     DefaultAction = "copy"
	DefaultActionDownload DefaultAction = "download"
)

type User struct {
	UserID        int64         `gorm:"column:user_id;primaryKey"`
	Username      string        `gorm:"column:username"`
	FirstName     string        `gorm:"column:first_name"`
	LastName      string        `gorm:"column:last_name"`
	LanguageCode  string        `gorm:"column:language_code"`
	DefaultAction DefaultAction `gorm:"column:default_action;default:'copy'"`
	IsActive      bool          `gorm:"column:is_active;default:true;index"`
	CreatedAt     time.Time     `gorm:"column:created_at;autoCreateTime"`
	LastSeenAt    time.Time     `gorm:"column:last_seen_at;autoCreateTime;index"`
}
