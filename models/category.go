package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string    `gorm:"unique;not null;"`
	Products []Product `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
