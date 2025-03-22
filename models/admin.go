package models

import "time"

type Admin struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Username         string     `gorm:"uniqueIndex;not null" json:"username"`
	Password         string     `gorm:"not null" json:"-"`	
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}