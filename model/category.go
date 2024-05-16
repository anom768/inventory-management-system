package model

import "gorm.io/gorm"

type Categories struct {
	gorm.Model
	Name          string `gorm:"column:name;not null" json:"name"`
	Specification string `gorm:"column:specification;not null" json:"specification"`
}
