package main

import (
	"fmt"
	"ktrhportal/database"
	"ktrhportal/routes"
	"ktrhportal/utilities"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

var AllowedOrigins = []string{
	"http://ktrh.tsavoerp.com",
	"http://localhost:8999",
	"http://localhost:9000",
	"http://localhost:5173",
	"http://102.37.201.86",
	"http://localhost:4001",
	"https://ktrh.go.ke",
	"http://ktrhportal.tsavoerp.com",
}

var AllowedHeaders = []string{
	"Authorization", "Accept", "Accept-Charset", "Accept-Language",
	"Accept-Encoding", "Origin", "Host", "User-Agent", "Content-Length",
	"Content-Type", " X-Authorization", "XMLHttpRequest", "Access-Control-Expose-Headers", " Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers", "Access-Control-Allow-Private-Network",
}

func main() {
	app := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = AllowedOrigins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "UPDATE"}
	config.AllowHeaders = AllowedHeaders
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	app.Use(cors.New(config))

	gin.ForceConsoleColor()
	database.Connect()
	routes.Setup(app)
	app.Run(fmt.Sprintf(":%s", utilities.GoDotEnvVariable("APP_PORT")))
}
