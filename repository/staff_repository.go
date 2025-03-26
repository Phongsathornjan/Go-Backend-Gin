package repository

import (
	"errors"
	"fmt"
	"os"
	"phongsathorn/go_backend_gin/models"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StaffRepository interface {
	CreateStaffID(username string, encrypted_password string, hospital_ID string) (status bool, error error)
	StaffLogin(username string, password string) (token string, status bool, error error)
}

type staffRepository struct {
	DB *gorm.DB
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{DB: db}
}

func (repo *staffRepository) CreateStaffID(username string, encrypted_password string, hospital_ID string) (status bool, error error) {
	var existingStaff models.Staff
	if err := repo.DB.Where("username = ?", username).First(&existingStaff).Error; err == nil {
		return false, fmt.Errorf("username already exists")
	} else if err != gorm.ErrRecordNotFound {
		return false, err
	}

	newStaff := models.Staff{
		Username:   username,
		Password:   encrypted_password,
		HospitalID: hospital_ID,
	}

	if err := repo.DB.Create(&newStaff).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (repo *staffRepository) StaffLogin(username string, password string) (token string, status bool, error error) {
	var staff models.Staff

	db_err := repo.DB.Where("username = ?", username).First(&staff).Error
	if db_err != nil {
		if errors.Is(db_err, gorm.ErrRecordNotFound) {
			return "", false, errors.New("username not found")
		}
		return "", false, db_err
	}

	bcrypt_err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(password))
	if bcrypt_err != nil {
		return "", false, errors.New("invalid password")
	}

	token, token_err := generateJWT(staff)
	if token_err != nil {
		return "", false, token_err
	}

	return token, true, nil
}

func generateJWT(staff models.Staff) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	fmt.Println(secretKey)
	claims := jwt.MapClaims{
		"staff_id":    staff.ID,
		"username":    staff.Username,
		"hospital_id": staff.HospitalID,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
