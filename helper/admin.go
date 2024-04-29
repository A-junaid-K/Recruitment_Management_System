package helper

import (
	"RMS_machine_task/db"
	"RMS_machine_task/domain/models"
	"log"
)

// Check the Admin aleady exists
func CheckAdminEmailExist(email string) (bool, error) {

	var count int64
	if err := db.DB.Table("admins").Where("email=?", email).Count(&count).Error; err != nil {
		log.Println("db email err: ", err)
		return true, err
	}
	return count > 0, nil
}

// Fetch the Admin by Email
func GetAdminByEmail(email string) (models.Admin, error) {
	var user models.Admin
	if err := db.DB.Table("admins").Where("email=?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
