package domain

import "time"

type Categories struct {
	ID            int       `gorm:"column:id;primary_key"`
	Name          string    `gorm:"column:name;not null" json:"name"`
	Specification string    `gorm:"column:specification;not null" json:"specification"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at;null" json:"deleted_at"`
}
