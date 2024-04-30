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

		user.POST("/add-education", middleware.ApplicantAuth, userHandler.AddEducation)
		user.POST("/add-experience", middleware.ApplicantAuth, userHandler.AddExperience)
	}

	router.POST("/extract-resume", userHandler.ExtractResume) // Third party API rate limit exceeded. processing

	job := router.Group("/user/job")
	{
		job.POST("/view", userHandler.ViewAllOpenings)
		job.POST("/apply/:job_id", middleware.ApplicantAuth, userHandler.Apply)
	}
}
