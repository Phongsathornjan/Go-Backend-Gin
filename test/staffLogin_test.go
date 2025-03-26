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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginTestCase struct {
	name           string
	prepareDB      func(*gorm.DB)
	requestBody    LoginRequest
	expectedStatus int
	expectedBody   map[string]interface{}
}

func SetupTestDatabase_staffLogin() (*gorm.DB, error) {
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

func prepareTestDatabase(t *testing.T) *gorm.DB {
	db, err := SetupTestDatabase_staffLogin()
	assert.NoError(t, err)
	db.Exec("DELETE FROM staffs")
	return db
}

func createTestStaff(db *gorm.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	staff := models.Staff{
		Username: username,
		Password: string(hashedPassword),
	}
	return db.Create(&staff).Error
}

func performLoginRequest(r *gin.Engine, requestBody LoginRequest) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(requestBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/staff/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func validateLoginResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody map[string]interface{}) {
	assert.Equal(t, expectedStatus, w.Code)

	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)

	for key, expectedValue := range expectedBody {
		assert.Equal(t, expectedValue, responseBody[key])
	}
}

func TestStaffLoginIntegration(t *testing.T) {
	testCases := []LoginTestCase{
		{
			name: "Successful Login",
			prepareDB: func(db *gorm.DB) {
				createTestStaff(db, "testuser", "correctpassword")
			},
			requestBody: LoginRequest{
				Username: "testuser",
				Password: "correctpassword",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "Login success",
			},
		},
		{
			name:      "Invalid Username",
			prepareDB: func(db *gorm.DB) {},
			requestBody: LoginRequest{
				Username: "nonexistentuser",
				Password: "somepassword",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"error": "username not found",
			},
		},
		{
			name: "Invalid Password",
			prepareDB: func(db *gorm.DB) {
				createTestStaff(db, "testuser", "correctpassword")
			},
			requestBody: LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"error": "invalid password",
			},
		},
		{
			name:      "Empty Credentials",
			prepareDB: func(db *gorm.DB) {},
			requestBody: LoginRequest{
				Username: "",
				Password: "",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"error": "username and password must be provided",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare database
			db := prepareTestDatabase(t)
			tc.prepareDB(db)

			// Setup router
			gin.SetMode(gin.TestMode)
			r := gin.Default()

			staffRepo := repository.NewStaffRepository(db)
			staffService := services.NewStaffService(staffRepo)
			staffController := controllers.NewStaffController(staffService)

			r.POST("/staff/login", staffController.StaffLogin)

			// Perform login request
			w := performLoginRequest(r, tc.requestBody)

			// Validate response
			validateLoginResponse(t, w, tc.expectedStatus, tc.expectedBody)
		})
	}
}
