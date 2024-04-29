package helper

import (
	"RMS_machine_task/db"
	"RMS_machine_task/domain/models"
	"errors"
	"log"
	"mime/multipart"
)

func CheckUserEmailExist(email string) (bool, error) {

	var count int64
	if err := db.DB.Table("users").Where("email=?", email).Count(&count).Error; err != nil {
		log.Println("db email err: ", err)
		return true, err
	}
	return count > 0, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := db.DB.Table("users").Where("email=?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func ValidateFileFormat(header *multipart.FileHeader) error {
	allowedFormats := map[string]bool{"application/pdf": true, "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true}

	contentType := header.Header.Get("Content-Type")
	if !allowedFormats[contentType] {
		return errors.New("invalid file format")
	}

	return nil
}
