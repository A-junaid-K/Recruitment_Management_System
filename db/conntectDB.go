package db

import (
	"RMS_machine_task/config"
	"RMS_machine_task/domain/models"

	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	Err error
)

func ConnectToDB(c config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort)
	DB, Err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if Err != nil {
		return nil, Err
	}

	DB.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.UserProfile{},

		&models.Admin{},
	)
	return DB, nil
}
