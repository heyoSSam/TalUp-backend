package main

import (
	"TalUp-backend/internal/auth"
	"TalUp-backend/internal/config/server"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := server.NewConfig()

	e := echo.New()

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
