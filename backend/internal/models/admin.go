package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type AdminToken struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	AdminID      uint      `gorm:"not null;index" json:"admin_id"`
	Admin        *Admin    `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	RefreshToken string    `gorm:"uniqueIndex;not null" json:"refresh_token"`
	JTI          string    `gorm:"index;not null" json:"jti"`
	ExpireAt     time.Time `gorm:"not null;index" json:"expire_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (t *AdminToken) IsExpired() bool {
	return time.Now().After(t.ExpireAt)
}

type AdminTokenBlacklist struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	JTI       string    `gorm:"uniqueIndex;not null" json:"jti"`
	ExpireAt  time.Time `gorm:"not null;index" json:"expire_at"`
	CreatedAt time.Time `json:"created_at"`
}
