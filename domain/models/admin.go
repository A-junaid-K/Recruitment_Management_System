package models

import "time"

type Admin struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string `json:"name,omitempty" validate:"min=3,max=20"`
	Email      string `json:"email,omitempty" validate:"email"`
	Password   string `json:"password,omitempty" validate:"min=6"`
	Otp        int
	UserType   string `json:"user_type" gorm:"NOT NULL"`
	Validate   bool   `json:"validate" gorm:"NOT NULL; default:false"`
	Created_at time.Time
}
