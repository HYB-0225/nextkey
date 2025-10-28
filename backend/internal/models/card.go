package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type Card struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CardKey    string         `gorm:"uniqueIndex;not null" json:"card_key"`
	ProjectID  uint           `gorm:"not null;index" json:"project_id"`
	Project    *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Activated  bool           `gorm:"default:false" json:"activated"`
	Duration   int            `gorm:"default:0" json:"duration"` // 秒
	ExpireAt   *time.Time     `json:"expire_at"`
	Note       string         `json:"note"`
	CardType   string         `gorm:"default:normal" json:"card_type"`
	CustomData string         `gorm:"type:text" json:"custom_data"` // JSON
	HWIDList   StringArray    `gorm:"type:text" json:"hwid_list"`
	IPList     StringArray    `gorm:"type:text" json:"ip_list"`
	MaxHWID    int            `gorm:"default:-1" json:"max_hwid"` // -1 无限制
	MaxIP      int            `gorm:"default:-1" json:"max_ip"`   // -1 无限制
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Card) IsExpired() bool {
	if c.ExpireAt == nil {
		return false
	}
	return time.Now().After(*c.ExpireAt)
}

func (c *Card) CanAddHWID() bool {
	if c.MaxHWID == -1 {
		return true
	}
	return len(c.HWIDList) < c.MaxHWID
}

func (c *Card) CanAddIP() bool {
	if c.MaxIP == -1 {
		return true
	}
	return len(c.IPList) < c.MaxIP
}
