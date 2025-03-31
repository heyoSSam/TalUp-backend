package main

import (
	_ "TalUp-backend/docs"
	"TalUp-backend/internal/auth"
	"TalUp-backend/internal/config/server"
	"TalUp-backend/internal/taskWord"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
)

// @title TalUp
// @version 1.0
// @description This is routes .
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {your JWT token}" to authenticate.

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := server.NewConfig()

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	authGroup := e.Group("/auth")
	taskWordGroup := e.Group("/taskWord")

	taskWord.RegisterRoutes(taskWordGroup)
	auth.RegisterRoutes(authGroup)

	if err := e.Start(":" + config.Port); err != nil {
		log.Fatal(err)
	}
}
