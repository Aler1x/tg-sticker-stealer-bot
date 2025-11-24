package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
	Packs *PackRepository
	Users *UserRepository
}

func New(dbPath string) (*DB, error) {
	conn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.AutoMigrate(&Pack{}, &User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	return &DB{
		DB:    conn,
		Packs: &PackRepository{db: conn},
		Users: &UserRepository{db: conn},
	}, nil
}

func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
