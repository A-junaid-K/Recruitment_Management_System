package routes

import (
	adminHandler "RMS_machine_task/controllers/adminHandler"
	jobHandler "RMS_machine_task/controllers/jobHandler"
	"RMS_machine_task/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {

	admin := router.Group("/admin")
	{
		admin.POST("/register", adminHandler.Register)
		admin.POST("/verify-otp", adminHandler.VerifyOtp)
		admin.POST("/login", adminHandler.AdminLogin)

		admin.GET("/applicants", middleware.AdminAuth, adminHandler.GetAllAplicants)
		admin.GET("/applicant/:applicant_id", middleware.AdminAuth, adminHandler.GetApplicantByID)
		admin.GET("/applicants/resume", middleware.AdminAuth, adminHandler.GetAllResume)
	}

	job := router.Group("/admin/job")
	{
		job.POST("/add-job", middleware.AdminAuth, jobHandler.AddJob)
		job.GET("/view/:id", middleware.AdminAuth, jobHandler.GetJobByParam)
		job.GET("/view/all", middleware.AdminAuth, jobHandler.GetAllJob)
	}

}
