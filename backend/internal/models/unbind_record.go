package models

import (
	"time"

	"gorm.io/gorm"
)

type UnbindRecord struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CardID       uint           `gorm:"not null;index" json:"card_id"`
	HWID         string         `gorm:"not null" json:"hwid"`
	UnbindAt     time.Time      `gorm:"not null" json:"unbind_at"`
	DeductedTime int            `gorm:"default:0" json:"deducted_time"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
