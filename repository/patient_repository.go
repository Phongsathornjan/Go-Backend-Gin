package repository

import (
	"phongsathorn/go_backend_gin/models"

	"gorm.io/gorm"
)

type PatientRepository interface {
	SearchPatient(patient models.Patient) ([]models.Patient, error)
	SearchPatientByID(id string) ([]models.Patient, error)
}

type patientRepository struct {
	DB *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{DB: db}
}

func (repo *patientRepository) SearchPatient(patient models.Patient) ([]models.Patient, error) {
	var patients []models.Patient
	tx := repo.DB

	if patient.NationalID != "" {
		tx = tx.Where("national_id = ?", patient.NationalID)
	}
	if patient.PassportID != "" {
		tx = tx.Where("passport_id = ?", patient.PassportID)
	}
	if patient.FirstNameTh != "" {
		tx = tx.Where("first_name_th ILIKE ?", "%"+patient.FirstNameTh+"%")
	}
	if patient.MiddleNameTh != "" {
		tx = tx.Where("middle_name_th ILIKE ?", "%"+patient.MiddleNameTh+"%")
	}
	if patient.LastNameTh != "" {
		tx = tx.Where("last_name_th ILIKE ?", "%"+patient.LastNameTh+"%")
	}
	if patient.DateOfBirth != "" {
		tx = tx.Where("date_of_birth = ?", patient.DateOfBirth)
	}
	if patient.PhoneNumber != "" {
		tx = tx.Where("phone_number = ?", patient.PhoneNumber)
	}
	if patient.Email != "" {
		tx = tx.Where("email ILIKE ?", "%"+patient.Email+"%")
	}

	tx = tx.Where("hospital_id = ?", patient.HospitalID)

	if err := tx.Find(&patients).Error; err != nil {
		return []models.Patient{}, err
	}

	return patients, nil
}

func (repo *patientRepository) SearchPatientByID(id string) ([]models.Patient, error) {
	var patients []models.Patient
	tx := repo.DB

	tx = tx.Where("national_id = ? OR passport_id = ?", id, id)

	if err := tx.Find(&patients).Error; err != nil {
		return []models.Patient{}, err
	}

	return patients, nil
}
