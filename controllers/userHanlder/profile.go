package userhandler

import (
	"RMS_machine_task/config"
	"RMS_machine_task/db"
	"RMS_machine_task/domain/models"
	"RMS_machine_task/domain/response"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
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
	if err := db.DB.Table("user_profiles").Where("applicant_id=?", user_id).Set("resume_file_address", result.Location).Error; err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "failed to store resume url in db", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	resp := response.SuccessResnpose{StatusCode: 200, Response: "Succefully Uploaded in S3 Bucket"}
	c.JSON(200, resp)
}
