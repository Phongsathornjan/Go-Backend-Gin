package controllers

import (
	"net/http"
	"phongsathorn/go_backend_gin/services"

	"github.com/gin-gonic/gin"
)

type StaffController struct {
	service services.StaffService
}

func NewStaffController(service services.StaffService) *StaffController {
	return &StaffController{service: service}
}

func (sc *StaffController) CreateStaffID(c *gin.Context) {

	var requestBody struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		Hospital_ID string `json:"hospital_id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := sc.service.CreateStaffID(requestBody.Username, requestBody.Password, requestBody.Hospital_ID)

	if !status {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Create Staff ID success"})

}

func (sc *StaffController) StaffLogin(c *gin.Context) {

	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, status, err := sc.service.StaffLogin(requestBody.Username, requestBody.Password)

	if !status {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Login success", "token": token})
}
