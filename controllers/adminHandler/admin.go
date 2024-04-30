package adminhandler

import (
	"RMS_machine_task/db"
	"RMS_machine_task/domain/models"
	"RMS_machine_task/domain/response"
	"RMS_machine_task/helper"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func Register(c *gin.Context) {
	var body models.AdminRegisterRequest

	if err := c.Bind(&body); err != nil {
		res := response.ErrResponse{StatusCode: 400, Response: "body bind error", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Validate
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "invalid input", Error: err.Error()}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	//Check User Exist
	exist, dberr := helper.CheckAdminEmailExist(body.Email)
	if dberr != nil {
		c.JSON(http.StatusNotFound, gin.H{"dberror": dberr})
		return
	}
	if exist {
		resp := response.ErrResponse{
			StatusCode: 400,
			Response:   "email already in use. Please login instead or use a different Email to sign up",
			Error:      "",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Hashing user entered password
	body.Password = helper.Hashpassword(body.Password)

	// Generating OTP
	otp := helper.GenerateOTP()

	err := db.DB.Create(&models.Admin{Name: body.Name, Email: body.Email, Password: body.Password, Otp: otp, UserType: "admin", Created_at: time.Now()}).Error
	if err != nil {
		log.Println("admin register error: ", err)
		c.JSON(400, err)
		return
	}

	// Sending generated otp to user email
	if err := helper.SendOtp(strconv.Itoa(otp), body.Email); err != nil {
		c.JSON(400, err)
		return
	}

	c.Redirect(303, "/verify-otp")
}

func VerifyOtp(c *gin.Context) {
	var otp models.VerifyOtp

	if err := c.Bind(&otp); err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "bind error", Error: err.Error()}
		c.JSON(http.StatusBadRequest, resp)
	}

	user, err := helper.GetAdminByEmail(otp.Email)
	if err != nil {
		c.JSON(400, err)
		return
	}
	if user.Email == "" {
		resp := response.ErrResponse{StatusCode: 400, Response: "Incorrect Email"}
		c.JSON(400, resp)
		return
	}

	if otp.Otp == user.Otp {
		//	Making admin validate = true
		if err := db.DB.Table("admins").Where("id=?", user.Id).Update("validate", true).Error; err != nil {
			c.JSON(400, err)
			return
		}

	} else {
		if err := db.DB.Table("admins").Where("id=?", user.Id).Delete(&models.User{}).Error; err != nil {
			c.JSON(500, err)
			return
		}
		resp := response.ErrResponse{StatusCode: 400, Response: "Incorrect OTP"}
		c.JSON(400, resp)
		return
	}

	resp := response.SuccessResnpose{StatusCode: 200, Response: "Succefully Registered"}
	c.JSON(200, resp)
}

func AdminLogin(c *gin.Context) {
	var loginBody models.LoginRequest

	if err := c.Bind(&loginBody); err != nil {
		res := response.ErrResponse{Response: "Binding Error", Error: err.Error(), StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	validator := validator.New()
	if err := validator.Struct(&loginBody); err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusBadRequest, Response: "Invalid Input", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	// Fetching the Admin details
	user, err := helper.GetAdminByEmail(loginBody.Email)
	if err != nil {
		res := response.ErrResponse{Response: "Wrong Email or Password", Error: err.Error(), StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := helper.ComapareHashPassword(user.Password, loginBody.Password); !err {
		res := response.ErrResponse{Response: "wrong Email or password", Error: "", StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	token, tokenerr := helper.CreateAdminAccessToken(&user, "admin")
	if tokenerr != nil {
		res := response.ErrResponse{Response: "failed to create access token", Error: tokenerr.Error(), StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	resp := models.LoginResopnse{StatusCode: http.StatusCreated, Token: token}
	c.JSON(201, resp)
}

type applicants struct {
	Name              string
	Email             string
	Headline          string
	ResumeFileAddress string
	Skills            string
	City              string
	State             string
	ZipCode           string
}

func GetAllAplicants(c *gin.Context) {
	var allusers []applicants

	if err := db.DB.Table("users").
		Select("users.name, users.email, user_profiles.headline, user_profiles.skills, user_profiles.resume_file_address, addresses.city, addresses.state, addresses.zip_code").
		Joins("JOIN user_profiles ON users.id = user_profiles.applicant_id").
		Joins("JOIN addresses ON users.id = addresses.applicant_id").
		Scan(&allusers).
		Error; err != nil {
		res := response.ErrResponse{Response: "Something wrong with findng applicants", Error: err.Error(), StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	c.JSON(200, allusers)
}

func GetApplicantByID(c *gin.Context) {
	// Recieve Parameter
	applicant_id := c.Param("applicant_id")

	// Using JOIN Method to get user and profile details
	var applicant applicants
	if err := db.DB.Table("users").Where("id=?", applicant_id).
		Select("users.name, users.email, user_profiles.headline, user_profiles.skills, user_profiles.resume_file_address, addresses.city, addresses.state, addresses.zip_code").
		Joins("JOIN user_profiles ON users.id = user_profiles.applicant_id").
		Joins("JOIN addresses ON users.id = addresses.applicant_id").
		Scan(&applicant).
		Error; err != nil {
		res := response.ErrResponse{Response: "Something wrong with findng applicants", Error: err.Error(), StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	c.JSONP(http.StatusOK, applicant)
}
