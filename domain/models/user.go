package models

import (
	"time"
)

type User struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string `json:"name,omitempty" validate:"min=3,max=20"`
	Email      string `json:"email,omitempty" validate:"email"`
	Password   string `json:"password,omitempty" validate:"min=6"`
	Otp        int
	UserType   string `json:"user_type" gorm:"NOT NULL"`
	Validate   bool   `json:"validate" gorm:"NOT NULL; default:false"`
	Created_at time.Time
	IsBlocked  bool `json:"isblocked" gorm:"default=false"`
}

type Address struct {
	Address_ID  uint `json:"address_id" gorm:"primaryKey;autoIncrement"`
	ApplicantID uint
	City        string `json:"city" gorm:"not null"`
	State       string `json:"state" gorm:"not null"`
	ZipCode     string `json:"zip_code" gorm:"not null"`
}

type UserProfile struct {
	ProfileID         int `json:"user_profile_id" gorm:"primaryKey;autoIncrement"`
	ApplicantID       uint
	Headline          string `json:"headline"`
	ResumeFileAddress string `json:"resume_file_address"`
	Skills            string `json:"skills"`

	Education           string
	EducationUrl        string
	EducationTimePeriod time.Time

	Experience           string
	CompanyUrl           string
	ExperienceTimePeriod time.Time

	Name  string `json:"name"`
	Email string `json:"email"`
	Phone int    `json:"phone"`
}
