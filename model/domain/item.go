package domain

import (
	"time"
)

type Items struct {
	ID            int       `gorm:"primaryKey;column:id;AUTO_INCREMENT"`
	Name          string    `gorm:"column:name;not null" json:"name"`
	CategoryID    int       `gorm:"column:category_id;not null" json:"category_id"`
	Quantity      int       `gorm:"column:quantity;not null" json:"quantity"`
	Price         float64   `gorm:"column:price;not null" json:"price"`
	Specification string    `gorm:"column:specification;type:text" json:"specification"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at"`
}
