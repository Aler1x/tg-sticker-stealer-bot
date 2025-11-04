package db

import (
	"fmt"

	"gorm.io/gorm"
)

// CreatePack creates a new pack in the database
func (r *Repository) CreatePack(pack *Pack) error {
	result := r.db.Create(pack)
	return result.Error
}

// GetPacksByUserID retrieves all packs for a specific user
func (r *Repository) GetPacksByUserID(userID int64) ([]Pack, error) {
	var packs []Pack
	result := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&packs)
	return packs, result.Error
}

// GetPackByID retrieves a specific pack by ID and user ID
func (r *Repository) GetPackByID(packID, userID int64) (*Pack, error) {
	var pack Pack
	result := r.db.Where("id = ? AND user_id = ?", packID, userID).First(&pack)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &pack, result.Error
}

// DeletePack deletes a pack by ID for a specific user
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

// GetPackCount returns the total number of packs in the database
func (r *Repository) GetPackCount() (int, error) {
	var count int64
	result := r.db.Model(&Pack{}).Count(&count)
	return int(count), result.Error
}
