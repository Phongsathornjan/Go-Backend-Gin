package models

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	PatientHN    string `gorm:"unique;not null"`
	FirstNameTh  string `gorm:"not null"`
	MiddleNameTh string
	LastNameTh   string `gorm:"not null"`
	FirstNameEn  string
	MiddleNameEn string
	LastNameEn   string
	DateOfBirth  string `gorm:"not null"`
	NationalID   string `gorm:"unique"`
	PassportID   string `gorm:"unique"`
	PhoneNumber  string `gorm:"index"`
	Email        string `gorm:"index"`
	Gender       string `gorm:"not null"`
	HospitalID   uint   `gorm:"not null"`
}

func (Patient) TableName() string {
	return "patient"
}
