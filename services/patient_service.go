package services

import (
	"fmt"
	"phongsathorn/go_backend_gin/models"
	"phongsathorn/go_backend_gin/repository"
)

type PatientService interface {
	SearchPatient(patient models.Patient) ([]models.Patient, error)
	SearchPatientByID(id string) ([]models.Patient, error)
}

type patientService struct {
	PatientRepo repository.PatientRepository
}

func NewPatientService(repo repository.PatientRepository) PatientService {
	return &patientService{PatientRepo: repo}
}

func (service *patientService) SearchPatient(patient models.Patient) ([]models.Patient, error) {

	if patient.HospitalID == 0 {
		return nil, fmt.Errorf("hospital id must be provided")
	}

	if patient.NationalID == "" && patient.PassportID == "" && patient.FirstNameTh == "" &&
		patient.MiddleNameTh == "" && patient.LastNameTh == "" && patient.DateOfBirth == "" &&
		patient.PhoneNumber == "" && patient.Email == "" {
		return nil, fmt.Errorf("provide at least one search criteria")
	}

	user, err := service.PatientRepo.SearchPatient(patient)
	if err != nil {
		fmt.Println("Error fetching user:", err)
		return []models.Patient{}, err
	}

	return user, nil
}

func (service *patientService) SearchPatientByID(id string) ([]models.Patient, error) {

	if id == "" {
		return nil, fmt.Errorf("ID must be provided")
	}

	user, err := service.PatientRepo.SearchPatientByID(id)
	if err != nil {
		fmt.Println("Error fetching user:", err)
		return []models.Patient{}, err
	}
	return user, nil
}
