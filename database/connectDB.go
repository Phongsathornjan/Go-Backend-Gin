package database

import (
	"log"
	"os"

	"phongsathorn/go_backend_gin/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_DSN")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = DB.AutoMigrate(&models.Staff{}, &models.Patient{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
