package models

import "gorm.io/gorm"

type Staff struct {
	gorm.Model
	Username   string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	HospitalID string `gorm:"not null"`
}

func (Staff) TableName() string {
	return "staff"
}
