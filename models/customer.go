package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name     string
	Password string
	Orders   []Order `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

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

type Category struct {
	gorm.Model
	Name     string    `gorm:"unique;not null;"`
	Products []Product `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type Salesman struct {
	gorm.Model
	Name     string
	Password string
	Orders   []Order `gorm:"foreignKey:SalesmanID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type Order struct {
	gorm.Model
	SalesmanID uint       `gorm:"not null;"`
	CustomerID uint       `gorm:"not null;"`
	Products   []*Product `gorm:"many2many:product_orders;"`
}
