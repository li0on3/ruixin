package models

import (
	"time"
)

type City struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CityID    int       `gorm:"uniqueIndex;not null" json:"city_id"`
	CityName  string    `gorm:"type:varchar(50);not null;index" json:"city_name"`
	Pinyin    string    `gorm:"type:varchar(10)" json:"pinyin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}