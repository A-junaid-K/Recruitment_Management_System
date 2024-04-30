package jobhandler

import (
	"RMS_machine_task/db"
	"RMS_machine_task/domain/models"
	"RMS_machine_task/domain/response"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func AddJob(c *gin.Context) {
	var body models.Job

	if err := c.Bind(&body); err != nil {
		res := response.ErrResponse{StatusCode: 400, Response: "body bind error", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Validate Admin Request
	validator := validator.New()
	if err := validator.Struct(&body); err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusBadRequest, Response: "Invalid Input", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	admin_id := c.GetInt("id")
	var admin models.Admin
	if err := db.DB.Table("admins").Where("id=?", admin_id).First(&admin).Error; err != nil {
		res := response.ErrResponse{StatusCode: 400, Response: "admin not found", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err := db.DB.Create(&models.Job{
		Title:       body.Title,
		Description: body.Description,
		CompanyName: body.CompanyName,
		PostedOn:    time.Now().Format("2006-01-02 15:04"),
		PostBy:      admin.Name,
	}).Error
	if err != nil {
		res := response.ErrResponse{StatusCode: 500, Response: "Failed to post job opening", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	resp := response.SuccessResnpose{StatusCode: 200, Response: "Succefully Posted the Job openings"}
	c.JSON(200, resp)
}

func GetJobByParam(c *gin.Context) {
	log.Println("before")
	id := c.Param("id")

	var filterjob models.Job
	err := db.DB.Table("jobs").Where("job_id=?", id).First(&filterjob).Error
	if err != nil {
		res := response.ErrResponse{StatusCode: 500, Response: "Failed to Get job opening by id", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}
	log.Println("after")

	resp := response.SuccessResnpose{StatusCode: 200, Response: filterjob}
	c.JSON(200, resp)
}

func GetAllJob(c *gin.Context) {
	type job struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Company     string `json:"company"`
		PostedOn    string `json:"posted_on"`
		PostBy      string `json:"post_by"`
	}
	var alljobs []job
	if err := db.DB.Table("jobs").Find(&alljobs).Error; err != nil {
		res := response.ErrResponse{StatusCode: 500, Response: "Error while feching job details", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	c.JSON(http.StatusOK, alljobs)
}
