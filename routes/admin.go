package routes

import (
	adminHandler "RMS_machine_task/controllers/adminHandler"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {

	r := router.Group("/admin")
	{
		r.POST("/register", adminHandler.Register)
		r.POST("/verify-otp", adminHandler.VerifyOtp)
		r.POST("/login", adminHandler.AdminLogin)

	}
}
