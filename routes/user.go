package routes

import (
	userHandler "RMS_machine_task/controllers/userHanlder"
	"RMS_machine_task/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {

	user := router.Group("/user")
	{
		user.POST("/signup", userHandler.UserSignup)
		user.POST("/verify-otp", userHandler.VerifyOtp)
		user.POST("/login", userHandler.UserLogin)

		user.POST("/upload-resume", middleware.ApplicantAuth, userHandler.UploadResume)
		user.PUT("/update/profile", middleware.ApplicantAuth, userHandler.UpdateUserProfile)
	}

	job := router.Group("/user/job")
	{
		job.POST("/view", userHandler.ViewAllOpenings)
		job.POST("/apply", middleware.ApplicantAuth, userHandler.Apply)
	}
}
