package helper

import (
	"RMS_machine_task/config"
	"RMS_machine_task/domain/models"

	"time"

	"github.com/golang-jwt/jwt"
)

func CreateAccessToken(user *models.User, user_type string) (accessToken string, err error) {
	cfg := config.GetConfig()
	exp := time.Now().Add(time.Hour * time.Duration(config.GetConfig().UserAccessTokenExpiryHour)).Unix()
	claims := &models.JwtCustomClaims{
		Email: user.Email,
		Id:    user.Id,
		User_type: user_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.UserAccessTokenSecret))
}

func CreateAdminAccessToken(admin *models.Admin, user_type string) (accessToken string, err error) {
	cfg := config.GetConfig()
	exp := time.Now().Add(time.Hour * time.Duration(config.GetConfig().AdminAccessTokenExpiryHour)).Unix()
	claims := &models.JwtCustomClaims{
		Email: admin.Email,
		Id:    admin.Id,
		User_type: user_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.AdminAccessTokenSecret))
}

