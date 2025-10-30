package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	UUID             string         `gorm:"uniqueIndex;not null" json:"uuid"`
	Name             string         `gorm:"not null" json:"name"`
	Mode             string         `gorm:"default:free" json:"mode"` // free/paid
	EnableHWID       bool           `gorm:"default:true" json:"enable_hwid"`
	EnableIP         bool           `gorm:"default:true" json:"enable_ip"`
	Version          string         `gorm:"default:1.0.0" json:"version"`
	UpdateURL        string         `json:"update_url"`
	TokenExpire      int            `gorm:"default:3600" json:"token_expire"`
	Description      string         `json:"description"`
	EnableUnbind     bool           `gorm:"default:false" json:"enable_unbind"`
	UnbindVerifyHWID bool           `gorm:"default:true" json:"unbind_verify_hwid"`
	UnbindDeductTime int            `gorm:"default:0" json:"unbind_deduct_time"`
	UnbindCooldown   int            `gorm:"default:86400" json:"unbind_cooldown"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
