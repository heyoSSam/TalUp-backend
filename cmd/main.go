package main

import (
	"TalUp-backend/internal/auth"
	"TalUp-backend/internal/config/server"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := server.NewConfig()

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	e := echo.New()

	authGroup := e.Group("/auth")
	auth.RegisterRoutes(authGroup)

	if err := e.Start(":" + config.Port); err != nil {
		log.Fatal(err)
	}
}
