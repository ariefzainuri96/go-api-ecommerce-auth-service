package entity

import (
	"gorm.io/gorm"
	"time"
)

// @Model
type BaseEntity struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	// DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}
