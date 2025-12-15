package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type PackRepository struct {
	db *gorm.DB
}

func (r *PackRepository) Create(pack *Pack) error {
	if err := r.db.Create(pack).Error; err != nil {
		return fmt.Errorf("failed to create pack: %w", err)
	}
	return nil
}

func (r *PackRepository) GetByUserID(userID int64) ([]Pack, error) {
	var packs []Pack
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&packs).Error; err != nil {
		return nil, fmt.Errorf("failed to query packs: %w", err)
	}
	return packs, nil
}

func (r *PackRepository) GetByID(packID, userID int64) (*Pack, error) {
	var pack Pack
	err := r.db.Where("id = ? AND user_id = ?", packID, userID).First(&pack).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pack: %w", err)
	}
	return &pack, nil
}

func (r *PackRepository) Delete(packID, userID int64) error {
	result := r.db.Where("id = ? AND user_id = ?", packID, userID).Delete(&Pack{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete pack: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("pack not found or not owned by user")
	}
	return nil
}

func (r *PackRepository) DeleteByPackName(packName string, userID int64) error {
	result := r.db.Where("pack_name = ? AND user_id = ?", packName, userID).Delete(&Pack{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete pack: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("pack not found or not owned by user")
	}
	return nil
}

func (r *PackRepository) GetByRelativeID(userID int64, relativeID int) (*Pack, error) {
	var pack Pack
	offset := relativeID - 1

	if offset < 0 {
		return nil, fmt.Errorf("invalid relative ID")
	}

	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(1).
		Offset(offset).
		First(&pack).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pack by relative ID: %w", err)
	}

	return &pack, nil
}

func (r *PackRepository) GetPaginated(userID int64, page, pageSize int) ([]Pack, int, error) {
	var packs []Pack
	var total int64

	if err := r.db.Model(&Pack{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count packs: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&packs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query packs: %w", err)
	}

	return packs, int(total), nil
}

func (r *PackRepository) Count() (int, error) {
	var count int64
	if err := r.db.Model(&Pack{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count packs: %w", err)
	}
	return int(count), nil
}
