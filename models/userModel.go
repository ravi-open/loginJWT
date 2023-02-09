package models

import "gorm.io/gorm"

type UserOpen struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}
