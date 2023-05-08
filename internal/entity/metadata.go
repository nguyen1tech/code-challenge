package entity

import (
	"time"

	"gorm.io/gorm"
)

type Metadata struct {
	ID        uint64         `gorm:"primaryKey" json:"id,string"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`

	Name     string `json:"name"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
	Location string `json:"location"`
}
