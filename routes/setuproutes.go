package routes

import (
	middleware "phongsathorn/go_backend_gin/Middleware"
	"phongsathorn/go_backend_gin/controllers"
	"phongsathorn/go_backend_gin/database"
	"phongsathorn/go_backend_gin/repository"
	"phongsathorn/go_backend_gin/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	patientRepo := repository.NewPatientRepository(database.DB)
	patientService := services.NewPatientService(patientRepo)
	patientController := controllers.NewPatientController(patientService)

	staffRepo := repository.NewStaffRepository(database.DB)
	staffService := services.NewStaffService(staffRepo)
	staffController := controllers.NewStaffController(staffService)

	patientRoutes := r.Group("/patient")
	{
		patientRoutes.POST("/search", middleware.TokenValidationMiddleware(), patientController.SearchPatient)
		patientRoutes.GET("/search/:id", patientController.SearchPatientByID)
	}

	stuffRoutes := r.Group("/staff")
	{
		stuffRoutes.POST("/create", middleware.TokenValidationMiddleware(), staffController.CreateStaffID)
		stuffRoutes.POST("/login", staffController.StaffLogin)
	}
}
