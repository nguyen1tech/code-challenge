package entity

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primaryKey" json:"id,string"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`

	Username string `json:"username"`
	Password string `json:"-"`
}

// GetID returns the user ID.
func (user User) GetID() string {
	return strconv.FormatUint(user.ID, 10)
}
