package models

import "github.com/golang-jwt/jwt"

type SignUpRequest struct {
	Name    string      `json:"name,omitempty" validate:"min=3,max=20"`
	Email   string      `json:"email,omitempty" validate:"email"`
	Profile UserProfile `json:"profile,omitempty"`
	Address Address     `json:"address,omitempty"`

	Password        string `json:"password,omitempty" validate:"min=6"`
	ConfirmPassword string `json:"confirmpassword" validate:"required,eqfield=Password"`
}

type VerifyOtp struct {
	Email string `json:"email,omitempty" validate:"email"`
	Otp   int    `json:"otp,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty" validate:"email"`
	Password string `json:"password,omitempty" validate:"min=6"`
}
type LoginResopnse struct {
	StatusCode int
	Token      string
}

type JwtCustomClaims struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	User_type string `json:"user_type"`
	jwt.StandardClaims
}

type JwtCustomRefreshClaim struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

type ChangePassword struct {
	Email           string `json:"email,omitempty" validate:"email"`
	OldPassword     string `json:"password,omitempty" validate:"min=6"`
	NewPassword     string `json:"newpassword,omitempty" validate:"min=6"`
	ConfirmPassword string `json:"confirmpassword" validate:"required,eqfield=Password"`
}
