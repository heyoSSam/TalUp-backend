package auth

import (
	"TalUp-backend/pkg/db"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Id             int       `json:"id"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Username       string    `json:"username"`
	CreatedAt      time.Time `json:"created_at"`
	LanguageLevel  string    `json:"language_level"`
	NativeLanguage string    `json:"native_language"`
	Xp             int       `json:"xp"`
	LastLogin      time.Time `json:"last_login"`
}

func RegisterRoutes(e *echo.Group) {
	e.POST("/login", Login)
	e.POST("/register", Register)
}

func Login(c echo.Context) error {
	return nil
}

func Register(c echo.Context) error {
	conn, err := db.GetDBConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	var user User

	if err := c.Bind(&user); err != nil {
		return err
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO users (id, email, password_hash, username, created_at, language_level, native_language, xp, last_login) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		user.Id, user.Email, hashedPassword, user.Username, user.CreatedAt, user.LanguageLevel, user.NativeLanguage, user.Xp, user.LastLogin)
	if err != nil {
		log.Fatal("Error inserting data:", err)
	}

	return c.JSON(http.StatusOK, user)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
