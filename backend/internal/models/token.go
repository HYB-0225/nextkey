package models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Token     string         `gorm:"uniqueIndex;not null" json:"token"`
	CardID    uint           `gorm:"not null;index" json:"card_id"`
	Card      *Card          `gorm:"foreignKey:CardID" json:"card,omitempty"`
	ProjectID uint           `gorm:"not null;index" json:"project_id"`
	ExpireAt  time.Time      `gorm:"not null" json:"expire_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpireAt)
}
