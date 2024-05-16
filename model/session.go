package model

import (
	"gorm.io/gorm"
	"time"
)

type Sessions struct {
	gorm.Model
	UserID    uint      `gorm:"column:user_id;not null" json:"user_id"`
	Token     string    `gorm:"column:token;unique;not null" json:"token"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime;<-:create" json:"created_at"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
}
