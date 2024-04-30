package userhandler

import (
	"RMS_machine_task/config"
	"RMS_machine_task/db"
	"RMS_machine_task/domain/models"
	"RMS_machine_task/domain/response"
	"RMS_machine_task/helper"
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func UploadResume(c *gin.Context) {
	user_id := c.GetInt("id")

	// Bind
	var body models.UserProfile
	if err := c.Bind(&body); err != nil {
		resp := response.ErrResponse{StatusCode: 500, Response: "Cannot Bind", Error: err.Error()}
		c.JSON(500, resp)
		return
	}

	resume, header, err := c.Request.FormFile("resume")
	if err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "Failed to get form file request", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	// validates the format of the uploaded file
	if err := helper.ValidateFileFormat(header); err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "Invalid file formal. Upload PDF or DOCX files", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	cfg := config.GetConfig()

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AwsRegion),
		Credentials: credentials.NewStaticCredentials(cfg.AwsAccessKey, cfg.AwsSecretAccessKey, ""),
	}))

	// filename_slice := strings.Split(header.Filename, ".")
	// ext := filename_slice[len(filename_slice)-1]
	uploader := s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cfg.AwsBucket),
		Key:    aws.String("resume/" + strconv.Itoa(user_id) + "." + header.Filename),
		ACL:    aws.String("public-read"),
		Body:   resume,
	})

	if err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "failed to upload image in s3 bucket", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	// store in db
	if err := db.DB.Table("user_profiles").Where("applicant_id=?", user_id).Update("resume_file_address", result.Location).Error; err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "failed to store resume url in db", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	resp := response.SuccessResnpose{StatusCode: 200, Response: "Succefully Uploaded in S3 Bucket"}
	c.JSON(200, resp)
}

func ViewAllOpenings(c *gin.Context) {

}

func Apply(c *gin.Context) {
	var body models.Apply

	// Bind
	if err := c.Bind(&body); err != nil {
		res := response.ErrResponse{Response: "Binding Error", Error: err.Error(), StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Validate User Request
	validator := validator.New()
	if err := validator.Struct(&body); err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusBadRequest, Response: "Invalid Input", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	// applicant_id from jwt
	applicant_id := c.GetInt("id")

	// Job details param
	job_id := c.Param("job_id")

	// Fetch Job details from DB
	var job models.Job
	err := db.DB.Table("jobs").Where("job_id=?", job_id).First(&job).Error
	if err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusBadRequest, Response: "Error while fetching job details", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	//Add Application to DB
	err = db.DB.Table("applies").Create(&models.Apply{
		ApplicantID:  uint(applicant_id),
		Name:         body.Name,
		Email:        body.Email,
		Phone:        body.Phone,
		Skills:       body.Skills,
		Education:    body.Education,
		Experience:   body.Experience,
		CurrentCTC:   body.CurrentCTC,
		ExpectedCTC:  body.ExpectedCTC,
		NoticePeriod: body.NoticePeriod,
	}).Error

	if err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusBadRequest, Response: "Filed to add application to db", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	// Add Application Count in Job database
	db.DB.Table("jobs").Where("job_id=?", job_id).Update("total_applicant", job.TotalApplicant+1)

	// Success response
	resp := response.SuccessResnpose{StatusCode: 200, Response: "Succefully Applied"}
	c.JSON(200, resp)
}

// --------   The third party API has exceeded daily/monthly rate limit  ----------
func ExtractResume(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Create a buffer to store the file content
	fileBytes := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileBytes, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}

	// Prepare multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating form file"})
		return
	}
	if _, err := io.Copy(part, fileBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing to form file"})
		return
	}
	writer.Close()

	// Create POST request to third-party API
	req, err := http.NewRequest("POST", "https://api.apilayer.com/resume_parser/upload", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("apikey", "gNiXyflsFu3WNYCz1ZCxdWDb7oQg1Nl1")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request"})
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}

	// Return response from third-party API
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
