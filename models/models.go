package models

import (
	"gorm.io/gorm"
)


type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
}

type Product struct {
	gorm.Model
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}