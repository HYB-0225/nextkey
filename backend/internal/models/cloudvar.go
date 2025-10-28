package models

import (
	"time"

	"gorm.io/gorm"
)

type CloudVar struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ProjectID uint           `gorm:"not null;index:idx_project_key" json:"project_id"`
	Project   *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Key       string         `gorm:"not null;index:idx_project_key" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
