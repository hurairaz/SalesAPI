package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	SalesmanID uint       `gorm:"not null;"`
	CustomerID uint       `gorm:"not null;"`
	Products   []*Product `gorm:"many2many:product_orders;"`
}
