package models

import (
	"time"
)

type Nonce struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Nonce     string    `gorm:"uniqueIndex;not null" json:"nonce"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
