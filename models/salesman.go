package models

import "gorm.io/gorm"

type Salesman struct {
	gorm.Model
	Name     string
	Password string
	Orders   []Order `gorm:"foreignKey:SalesmanID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
