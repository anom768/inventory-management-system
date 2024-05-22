package domain

import "time"

type Categories struct {
	ID        int       `gorm:"primaryKey;column:id;AUTO_INCREMENT"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at;null"`
}
