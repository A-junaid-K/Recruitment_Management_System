package userhandler

import (
	"RMS_machine_task/db"
	"RMS_machine_task/domain/models"
	"RMS_machine_task/domain/response"
	"RMS_machine_task/helper"
	"log"
	"time"

	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func UserSignup(c *gin.Context) {
	var body models.SignUpRequest

	if err := c.Bind(&body); err != nil {
		res := response.ErrResponse{StatusCode: 400, Response: "body bind error", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	log.Println("log : ", body)

	// Validating User Entered data
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "Invalid input. Try again", Error: err.Error()}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	//Check if User already exists
	exist, dberr := helper.CheckUserEmailExist(body.Email)
	if dberr != nil {
		c.JSON(http.StatusNotFound, gin.H{"dberror": dberr})
		return
	}
	if exist {
		resp := response.ErrResponse{
			StatusCode: 400,
			Response:   "Email already in use. Please login instead or use a different Email to sign up",
			Error:      "",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Hashing user entered password
	body.Password = helper.Hashpassword(body.Password)

	// Generating OTP
	otp := helper.GenerateOTP()

	// Create User
	err := db.DB.Create(&models.User{
		Name:       body.Name,
		Email:      body.Email,
		Password:   body.Password,
		Otp:        otp,
		UserType:   "applicant",
		Created_at: time.Now(),
	}).Error

	if err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "Failed to create user", Error: err.Error()}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	signeduser, err := helper.GetUserByEmail(body.Email)
	if err != nil {
		log.Println("Failed to User Details by email : ", err)
	}

	// Upload User Profile in DB
	userProfile := &models.UserProfile{
		ApplicantID: uint(signeduser.Id),
		Headline:    body.Profile.Headline,
		Name:        body.Name,
		Email:       body.Email,
	}
	if err := db.DB.Create(userProfile).Error; err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "Failed to Create user profile", Error: err.Error()}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	log.Println("created user profile")

	// Upload User Address in DB
	userAddress := &models.Address{
		ApplicantID: uint(signeduser.Id),
		City:        body.Address.City,
		State:       body.Address.State,
		ZipCode:     body.Address.ZipCode,
	}

	if err := db.DB.Create(userAddress).Error; err != nil {
		resp := response.ErrResponse{StatusCode: 400, Response: "Failed to Create user address", Error: err.Error()}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Sending OTP to user email
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

	user, err := helper.GetUserByEmail(otp.Email)
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
		//	Making user validate = true
		if err := db.DB.Table("users").Where("id=?", user.Id).Update("validate", true).Error; err != nil {
			c.JSON(400, err)
			return
		}

	} else {
		// Invalid OTP
		if err := db.DB.Table("users").Where("id=?", user.Id).Delete(&models.User{}).Error; err != nil {
			c.JSON(500, err)
			return
		}
		resp := response.ErrResponse{StatusCode: 400, Response: "Incorrect OTP"}
		c.JSON(400, resp)
		return
	}

	resp := response.SuccessResnpose{StatusCode: 200, Response: "Succefully Signed Up"}
	c.JSON(200, resp)
}

func UserLogin(c *gin.Context) {
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

	// Fetching the User details
	user, err := helper.GetUserByEmail(loginBody.Email)
	if err != nil {
		res := response.ErrResponse{Response: "Wrong Email or Password", Error: err.Error(), StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	//Cheking User blocked or not
	if user.IsBlocked {
		res := response.ErrResponse{Response: "User blocked by Admin", Error: "", StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := helper.ComapareHashPassword(user.Password, loginBody.Password); !err {
		res := response.ErrResponse{Response: "wrong Email or password", Error: "", StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	token, tokenerr := helper.CreateAccessToken(&user, "applicant")
	if tokenerr != nil {
		res := response.ErrResponse{Response: "failed to create access token", Error: tokenerr.Error(), StatusCode: 400}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	resp := models.LoginResopnse{StatusCode: http.StatusCreated, Token: token}
	c.JSON(201, resp)
}

func UpdateUserProfile(c *gin.Context) {
	var body models.UserProfile

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

	// Applicant from JWT
	applicant_id := c.GetInt("id")

	// Retrieve existing user profile
	var existingProfile models.UserProfile
	if err := db.DB.Table("user_profiles").Where("applicant_id=?", applicant_id).First(&existingProfile).Error; err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusInternalServerError, Response: "Error while retrieving user profile", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Merge existing data with new data
	switch {
	case body.Headline == "":
		body.Headline = existingProfile.Headline
	case len(body.Skills) == 0:
		body.Skills = existingProfile.Skills
	case body.Phone == 0:
		body.Phone = existingProfile.Phone
	}

	// Update User Profile
	err := db.DB.Table("user_profiles").Where("applicant_id=?", applicant_id).Updates(map[string]interface{}{
		"headline": body.Headline,
		"skills":   body.Skills,
		"phone":    body.Phone,
	}).Error

	if err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusBadRequest, Response: "Erroe while updating user profile", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	resp := response.SuccessResnpose{StatusCode: 200, Response: "User profile successfully updated"}
	c.JSON(200, resp)
}

// Add Applicant Education
func AddEducation(c *gin.Context) {
	var body models.UserProfile

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

	// Applicant from JWT
	applicant_id := c.GetInt("id")

	// Retrieve existing user profile
	var existingProfile models.UserProfile
	if err := db.DB.Table("user_profiles").Where("applicant_id=?", applicant_id).First(&existingProfile).Error; err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusInternalServerError, Response: "Error while retrieving user profile", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Merge existing data with new data
	switch {
	case body.Education == "":
		body.Education = existingProfile.Education

	}

	// Update User Profile
	err := db.DB.Table("user_profiles").Where("applicant_id=?", applicant_id).Updates(map[string]interface{}{
		"education":             body.Education,
		"education_url":         body.EducationUrl,
		"education_time_period": body.EducationTimePeriod,
	}).Error

	if err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusBadRequest, Response: "Erroe while add Education", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	resp := response.SuccessResnpose{StatusCode: 200, Response: "Successfully added Applicant Education"}
	c.JSON(200, resp)
}

func AddExperience(c *gin.Context) {
	var body models.UserProfile

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

	// Applicant from JWT
	applicant_id := c.GetInt("id")

	// Retrieve existing user profile
	var existingProfile models.UserProfile
	if err := db.DB.Table("user_profiles").Where("applicant_id=?", applicant_id).First(&existingProfile).Error; err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusInternalServerError, Response: "Error while add Experience", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Merge existing data with new data
	switch {
	case body.Experience == "":
		body.Experience = existingProfile.Experience
	}

	// Update User Profile
	err := db.DB.Table("user_profiles").Where("applicant_id=?", applicant_id).Updates(map[string]interface{}{
		"experience":             body.Experience,
		"company_url":            body.CompanyUrl,
		"experience_time_period": body.ExperienceTimePeriod,
	}).Error

	if err != nil {
		resp := response.ErrResponse{StatusCode: http.StatusBadRequest, Response: "Erroe while updating user profile", Error: err.Error()}
		c.JSON(400, resp)
		return
	}

	resp := response.SuccessResnpose{StatusCode: 200, Response: "Successfully added Applicant Experience"}
	c.JSON(200, resp)
}
