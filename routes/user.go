package routes

import (
	userHandler "RMS_machine_task/controllers/userHanlder"
	"RMS_machine_task/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {

	r := router.Group("/user")
	{
		r.POST("/signup", userHandler.UserSignup)
		r.POST("/verify-otp", userHandler.VerifyOtp)
		r.POST("/login", userHandler.UserLogin)
		r.POST("/upload-resume", middleware.ApplicantAuth, userHandler.UploadResume)

		// r.POST("/profile", middlware.UserAuth, userHandler.UserProfile)
	}
}
