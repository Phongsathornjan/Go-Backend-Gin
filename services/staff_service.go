package services

import (
	"fmt"
	"phongsathorn/go_backend_gin/repository"

	"golang.org/x/crypto/bcrypt"
)

type StaffService interface {
	CreateStaffID(username string, password string, hospital_ID string) (status bool, error error)
	StaffLogin(username string, password string) (token string, status bool, error error)
}

type staffService struct {
	StaffRepo repository.StaffRepository
}

func NewStaffService(repo repository.StaffRepository) StaffService {
	return &staffService{StaffRepo: repo}
}

func (service *staffService) CreateStaffID(username string, password string, hospital_id string) (status bool, error error) {

	if username == "" || password == "" || hospital_id == "" {
		return false, fmt.Errorf("username, password, and hospital id must be provided")
	}

	hashedPassword, encrypt_err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if encrypt_err != nil {
		return false, fmt.Errorf("failed to encrypt password: %v", encrypt_err)
	}

	status, err := service.StaffRepo.CreateStaffID(username, string(hashedPassword), hospital_id)

	if !status {
		return false, err
	}

	return true, nil
}

func (service *staffService) StaffLogin(username string, password string) (token string, status bool, error error) {

	if username == "" || password == "" {
		return "", false, fmt.Errorf("username and password must be provided")
	}

	token, status, err := service.StaffRepo.StaffLogin(username, password)

	if !status {
		return "", false, err
	}

	return token, true, nil
}
