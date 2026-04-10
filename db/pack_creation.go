package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PackCreation struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	UserID   int64     `gorm:"column:user_id;not null;index"`
	PackLink string    `gorm:"column:pack_link;not null"`
	PackType string    `gorm:"column:pack_type;not null;index"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index"`
}

type PackCreationRepository struct {
	db *gorm.DB
}

func (r *PackCreationRepository) Record(userID int64, packLink, packType string) error {
	row := &PackCreation{
		ID:       uuid.New(),
		UserID:   userID,
		PackLink: packLink,
		PackType: packType,
	}
	if err := r.db.Create(row).Error; err != nil {
		return fmt.Errorf("failed to record pack creation: %w", err)
	}
	return nil
}

type PackTypeAggregate struct {
	PackType string `gorm:"column:pack_type"`
	Count    int64  `gorm:"column:cnt"`
}

func (r *PackCreationRepository) TotalCount() (int64, error) {
	var n int64
	if err := r.db.Model(&PackCreation{}).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("failed to count pack creations: %w", err)
	}
	return n, nil
}

func (r *PackCreationRepository) CountByPackType() ([]PackTypeAggregate, error) {
	var rows []PackTypeAggregate
	err := r.db.Model(&PackCreation{}).
		Select("pack_type, COUNT(*) AS cnt").
		Group("pack_type").
		Order("pack_type").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate pack types: %w", err)
	}
	return rows, nil
}
