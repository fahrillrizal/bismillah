package models

import "time"

type Category struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"notnull" json:"name"`
	Order int    `json:"order"`
	Links []Link `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"links,omitempty"`
}

type Link struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"not null" json:"title"`
	URL        string    `gorm:"not null" json:"url"`
	ImageURL   string    `json:"image_url"`
	Price      int64     `json:"price"`
	PriceStr   string    `json:"price_str"`
	Order      int       `json:"order"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CategoryID uint      `json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}