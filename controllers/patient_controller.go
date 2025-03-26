package controllers

import (
	"net/http"
	"phongsathorn/go_backend_gin/models"
	"phongsathorn/go_backend_gin/services"

	"github.com/gin-gonic/gin"
)

type PatientController struct {
	service services.PatientService
}

func NewPatientController(service services.PatientService) *PatientController {
	return &PatientController{service: service}
}

func (pc *PatientController) SearchPatient(c *gin.Context) {
	var patients []models.Patient

	var requestBody struct {
		NationalID  string `json:"national_id"`
		PassportID  string `json:"passport_id"`
		FirstName   string `json:"first_name"`
		MiddleName  string `json:"middle_name"`
		LastName    string `json:"last_name"`
		DateOfBirth string `json:"date_of_birth"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
		HospitalID  int    `json:"hospital_id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	patient := models.Patient{
		NationalID:   requestBody.NationalID,
		PassportID:   requestBody.PassportID,
		FirstNameTh:  requestBody.FirstName,
		MiddleNameTh: requestBody.MiddleName,
		LastNameTh:   requestBody.LastName,
		DateOfBirth:  requestBody.DateOfBirth,
		PhoneNumber:  requestBody.PhoneNumber,
		Email:        requestBody.Email,
		HospitalID:   uint(requestBody.HospitalID),
	}

	patients, err := pc.service.SearchPatient(patient)

	if patients == nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patients": patients})
}

func (pc *PatientController) SearchPatientByID(c *gin.Context) {
	var patients []models.Patient

	id := c.Param("id")
	patients, err := pc.service.SearchPatientByID(id)

	if patients == nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patients": patients})
}
