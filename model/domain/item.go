package domain

import "gorm.io/gorm"

type Items struct {
	gorm.Model
	Name        string  `gorm:"column:name;not null" json:"name"`
	CategoryID  int     `gorm:"column:category_id;not null" json:"category_id"`
	Quantity    int     `gorm:"column:quantity;not null" json:"quantity"`
	Price       float64 `gorm:"column:price;not null" json:"price"`
	Description string  `gorm:"column:description;type:text" json:"description"`
}
