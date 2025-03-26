package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"phongsathorn/go_backend_gin/controllers"
	"phongsathorn/go_backend_gin/models"
	"phongsathorn/go_backend_gin/repository"
	"phongsathorn/go_backend_gin/services"
)

type CreateStaffRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Hospital_ID string `json:"hospital_id"`
}

type CreateStaffTestCase struct {
	name           string
	prepareDB      func(*gorm.DB)
	requestBody    CreateStaffRequest
	expectedStatus int
	expectedBody   map[string]interface{}
}

func performCreateStaffRequest(r *gin.Engine, requestBody CreateStaffRequest) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(requestBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func validateCreateStaffResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody map[string]interface{}) {
	assert.Equal(t, expectedStatus, w.Code)

	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)

	for key, expectedValue := range expectedBody {
		assert.Equal(t, expectedValue, responseBody[key])
	}
}

func SetupTestDatabase_staffCreate() (*gorm.DB, error) {
	dsn := "host=localhost user=admin password=123456 dbname=mydb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Staff{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestCreateStaffIDIntegration(t *testing.T) {
	// Setup database
	db, err := SetupTestDatabase_staffCreate()
	assert.NoError(t, err)
	defer db.Exec("DELETE FROM staffs")

	testCases := []CreateStaffTestCase{
		{
			name: "Successful Staff ID Creation",
			prepareDB: func(db *gorm.DB) {
				// No preparation needed
			},
			requestBody: CreateStaffRequest{
				Username:    "newuser",
				Password:    "strongpassword",
				Hospital_ID: "HOSP001",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "Create Staff ID success",
			},
		},
		{
			name: "Duplicate Username",
			prepareDB: func(db *gorm.DB) {
				// Pre-create a user with the same username
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("existingpassword"), bcrypt.DefaultCost)
				existingStaff := models.Staff{
					Username:   "existinguser",
					Password:   string(hashedPassword),
					HospitalID: "HOSP002",
				}
				db.Create(&existingStaff)
			},
			requestBody: CreateStaffRequest{
				Username:    "existinguser",
				Password:    "newpassword",
				Hospital_ID: "HOSP003",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"error": "username already exists",
			},
		},
		{
			name: "Empty Username",
			prepareDB: func(db *gorm.DB) {
				// No preparation needed
			},
			requestBody: CreateStaffRequest{
				Username:    "",
				Password:    "strongpassword",
				Hospital_ID: "HOSP001",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"error": "username, password, and hospital id must be provided",
			},
		},
		{
			name: "Empty Password",
			prepareDB: func(db *gorm.DB) {
				// No preparation needed
			},
			requestBody: CreateStaffRequest{
				Username:    "newuser",
				Password:    "",
				Hospital_ID: "HOSP001",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"error": "username, password, and hospital id must be provided",
			},
		},
		{
			name: "Empty Hospital ID",
			prepareDB: func(db *gorm.DB) {
				// No preparation needed
			},
			requestBody: CreateStaffRequest{
				Username:    "newuser",
				Password:    "strongpassword",
				Hospital_ID: "",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"error": "username, password, and hospital id must be provided",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear database
			db.Exec("DELETE FROM staffs")

			// Prepare database
			tc.prepareDB(db)

			// Setup router
			gin.SetMode(gin.TestMode)
			r := gin.Default()

			staffRepo := repository.NewStaffRepository(db)
			staffService := services.NewStaffService(staffRepo)
			staffController := controllers.NewStaffController(staffService)

			r.POST("/staff/create", staffController.CreateStaffID)

			// Perform create staff request
			w := performCreateStaffRequest(r, tc.requestBody)

			// Validate response
			validateCreateStaffResponse(t, w, tc.expectedStatus, tc.expectedBody)

			// Additional verification for successful staff creation
			if tc.name == "Successful Staff ID Creation" {
				var staff models.Staff
				result := db.Where("username = ?", tc.requestBody.Username).First(&staff)
				assert.NoError(t, result.Error)
				assert.Equal(t, tc.requestBody.Username, staff.Username)
				assert.Equal(t, tc.requestBody.Hospital_ID, staff.HospitalID)

				// Verify password
				err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(tc.requestBody.Password))
				assert.NoError(t, err)
			}
		})
	}
}
