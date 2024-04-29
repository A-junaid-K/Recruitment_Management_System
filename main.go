package main

import (
	"RMS_machine_task/config"
	"RMS_machine_task/db"
	"RMS_machine_task/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.InitConfig()
	db.ConnectToDB(cfg)

	router := gin.Default()

	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	router.Run(":4000")
}
