package model

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Username string `gorm:"column:username;unique" json:"username"`
	Password string `gorm:"column:password" json:"password"`
	Role     string `gorm:"column:role" json:"role"`
}
