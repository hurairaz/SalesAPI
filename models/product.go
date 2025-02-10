package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Price       int
	Stock       int
	Image       string
	Description string
	Orders      []*Order `gorm:"many2many:product_orders"`
	CategoryID  uint
}
