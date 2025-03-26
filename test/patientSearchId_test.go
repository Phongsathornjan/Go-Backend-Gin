package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"phongsathorn/go_backend_gin/models"
	"phongsathorn/go_backend_gin/repository"
	"phongsathorn/go_backend_gin/services"
)

func setupTestDatabase_Patient() (*gorm.DB, error) {
	dsn := "host=localhost user=admin password=123456 dbname=mydb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Drop table if exists and recreate
	err = db.Migrator().DropTable(&models.Patient{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate the Patient model
	err = db.AutoMigrate(&models.Patient{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func clearPatientTable(db *gorm.DB) error {
	return db.Where("1 = 1").Delete(&models.Patient{}).Error
}

func insertTestPatients(db *gorm.DB) error {
	patients := []models.Patient{
		{
			PatientHN:   "HN001",
			FirstNameTh: "สมชาย",
			LastNameTh:  "ใจดี",
			FirstNameEn: "Somchai",
			LastNameEn:  "Jaidee",
			NationalID:  "1234567890123",
			PassportID:  "",
			DateOfBirth: "1990-01-01",
			Gender:      "Male",
			HospitalID:  1,
			PhoneNumber: "0812345678",
			Email:       "somchai@example.com",
		},
		{
			PatientHN:   "HN002",
			FirstNameTh: "สมหญิง",
			LastNameTh:  "แก้วใส",
			FirstNameEn: "Somying",
			LastNameEn:  "Kaewsai",
			NationalID:  "",
			PassportID:  "AB1234567",
			DateOfBirth: "1995-05-05",
			Gender:      "Female",
			HospitalID:  1,
			PhoneNumber: "0987654321",
			Email:       "somying@example.com",
		},
	}

	return db.Create(&patients).Error
}

func TestPatientSearchByID(t *testing.T) {
	// Setup test database
	db, err := setupTestDatabase_Patient()
	assert.NoError(t, err, "Failed to setup test database")

	// Clear existing data and insert test patients
	err = clearPatientTable(db)
	assert.NoError(t, err, "Failed to clear patient table")

	err = insertTestPatients(db)
	assert.NoError(t, err, "Failed to insert test patients")

	// Create repository and service
	patientRepo := repository.NewPatientRepository(db)
	patientService := services.NewPatientService(patientRepo)

	testCases := []struct {
		name            string
		inputID         string
		expectedCount   int
		expectedError   bool
		expectedPatient *models.Patient
	}{
		{
			name:          "Search by National ID",
			inputID:       "1234567890123",
			expectedCount: 1,
			expectedError: false,
			expectedPatient: &models.Patient{
				PatientHN:   "HN001",
				FirstNameTh: "สมชาย",
				LastNameTh:  "ใจดี",
				FirstNameEn: "Somchai",
				LastNameEn:  "Jaidee",
				NationalID:  "1234567890123",
			},
		},
		{
			name:          "Search by Passport ID",
			inputID:       "AB1234567",
			expectedCount: 1,
			expectedError: false,
			expectedPatient: &models.Patient{
				PatientHN:   "HN002",
				FirstNameTh: "สมหญิง",
				LastNameTh:  "แก้วใส",
				FirstNameEn: "Somying",
				LastNameEn:  "Kaewsai",
				PassportID:  "AB1234567",
			},
		},
		{
			name:            "Non-existent ID",
			inputID:         "999999999",
			expectedCount:   0,
			expectedError:   false,
			expectedPatient: nil,
		},
		{
			name:            "Empty ID",
			inputID:         "",
			expectedCount:   0,
			expectedError:   true,
			expectedPatient: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute search
			patients, err := patientService.SearchPatientByID(tc.inputID)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, patients)
			} else {
				assert.NoError(t, err)
				assert.Len(t, patients, tc.expectedCount)

				if tc.expectedCount > 0 && tc.expectedPatient != nil {
					patient := patients[0]
					assert.Equal(t, tc.expectedPatient.PatientHN, patient.PatientHN)
					assert.Equal(t, tc.expectedPatient.FirstNameTh, patient.FirstNameTh)
					assert.Equal(t, tc.expectedPatient.LastNameTh, patient.LastNameTh)
					assert.Equal(t, tc.expectedPatient.FirstNameEn, patient.FirstNameEn)
					assert.Equal(t, tc.expectedPatient.LastNameEn, patient.LastNameEn)

					// Check which ID was used for search
					if tc.inputID == tc.expectedPatient.NationalID {
						assert.Equal(t, tc.expectedPatient.NationalID, patient.NationalID)
					} else if tc.inputID == tc.expectedPatient.PassportID {
						assert.Equal(t, tc.expectedPatient.PassportID, patient.PassportID)
					}
				}
			}
		})
	}
}
