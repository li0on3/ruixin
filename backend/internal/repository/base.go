package repository

import (
	"errors"
	"gorm.io/gorm"
)

var (
	// ErrNoAvailableCard is returned when no available card is found
	ErrNoAvailableCard = errors.New("no available card found")
)

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}