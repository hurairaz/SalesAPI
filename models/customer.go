package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name     string
	Password string
	Orders   []Order `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
