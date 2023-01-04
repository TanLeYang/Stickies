package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection() (*gorm.DB, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	return gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
}
