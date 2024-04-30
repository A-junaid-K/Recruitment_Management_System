package middleware

import (
	"RMS_machine_task/config"
	"RMS_machine_task/domain/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AdminAuth(c *gin.Context) {
	cfg := config.GetConfig()
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		err := response.ErrResponse{StatusCode: http.StatusUnauthorized, Response: "Please provide your admin token", Error: "Empty Token"}
		c.JSON(404, err)
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.AdminAccessTokenSecret), nil
	})

	if err != nil || !token.Valid {
		resp := response.ErrResponse{StatusCode: 401, Response: "Cannot parse authorization token", Error: err.Error()}
		c.JSON(401, resp)
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		resp := response.ErrResponse{StatusCode: 401, Response: "Invalid Authorization token"}
		c.JSON(http.StatusUnauthorized, resp)
		c.Abort()
		return
	}

	role, ok := claims["user_type"].(string)
	if !ok || role != "admin" {
		resp := response.ErrResponse{StatusCode: 403, Response: "UnAuthorized Access"}
		c.JSON(http.StatusForbidden, resp)
		c.Abort()
		return
	}

	id, ok := claims["id"].(float64)
	if !ok || id == 0 {
		resp := response.ErrResponse{StatusCode: 403, Response: "Something wrong in token"}
		c.JSON(http.StatusForbidden, resp)
		c.Abort()
		return
	}
	uid := int(id)
	c.Set("id", uid)

}
