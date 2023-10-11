package main

import (
	"fmt"
	"ktrhportal/database"
	"ktrhportal/routes"
	"ktrhportal/utilities"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {
	app := gin.Default()
	gin.ForceConsoleColor()
	database.Connect()
	routes.Setup(app)
	app.Run(fmt.Sprintf(":%s", utilities.GoDotEnvVariable("APP_PORT")))
}
