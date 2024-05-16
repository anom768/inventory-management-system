package model

import (
	"gorm.io/gorm"
	"time"
)

type Sessions struct {
	gorm.Model
	Username  string    `gorm:"column:username;not null" json:"username"`
	Token     string    `gorm:"column:token;unique;not null" json:"token"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
}
