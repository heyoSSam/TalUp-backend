package main

import (
	_ "TalUp-backend/docs"
	"TalUp-backend/internal/auth"
	"TalUp-backend/internal/config/server"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := server.NewConfig()

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	authGroup := e.Group("/auth")
	auth.RegisterRoutes(authGroup)

	//Middleware for routes with token access
	//homeGroup.Use(echojwt.WithConfig(echojwt.Config{
	//	SigningKey: []byte(os.Getenv("JWT_SECRET")),
	//}))

	if err := e.Start(":" + config.Port); err != nil {
		log.Fatal(err)
	}
}
