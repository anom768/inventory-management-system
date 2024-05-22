package domain

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	FullName string `gorm:"column:full_name" json:"full_name"`
	Username string `gorm:"column:username;unique" json:"username"`
	Password string `gorm:"column:password"`
	Role     string `gorm:"column:role" json:"role"`
}
