package db

import (
	"fmt"

	"github.com/TanLeYang/stickies/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection(dbConf config.DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		dbConf.DbHost,
		dbConf.DbPort,
		dbConf.DbUser,
		dbConf.DbName,
		dbConf.DbPassword,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
