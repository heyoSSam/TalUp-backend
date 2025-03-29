package auth

import (
	"TalUp-backend/pkg/db"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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

type JwtCustomClaims struct {
	Username string `json:"username"`
	jwt.Claims
}

type LoginRequest struct {
	Email    string `json:"email" example:"cs@example.com"`
	Password string `json:"password" example:"1234"`
}

type RegRequest struct {
	Email          string `json:"email" example:"blablabla@gmail.com"`
	LanguageLevel  string `json:"language_level" example:"1"`
	NativeLanguage string `json:"native_language" example:"russian"`
	Password       string `json:"password" example:"1234"`
	Username       string `json:"username" example:"bla"`
}

func RegisterRoutes(e *echo.Group) {
	e.POST("/login", Login)
	e.POST("/register", Register)
}

// Login handles user authentication
// @Summary      Login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "User Credentials"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/login [post]
func Login(c echo.Context) error {
	var user User

	if err := c.Bind(&user); err != nil {
		return err
	}

	password, err := getUserPassword(user.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	if !passwordComparison(user.Password, password) {
		return c.JSON(http.StatusUnauthorized, "password is incorrect")
	} else {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 1).Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to generate token",
			})
		}

		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = tokenString
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.Path = "/"
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]string{
			"token":   tokenString,
			"message": "login success",
		})
	}
}

// Register handles user signup
// @Summary      Register
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegRequest true "User Data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/register [post]
func Register(c echo.Context) error {
	conn, err := db.GetDBConnection()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
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

	_, err = conn.Exec(context.Background(), "INSERT INTO users (email, password_hash, username, created_at, language_level, native_language, xp, last_login) VALUES ($1, $2, $3, $4, $5, $6, 0, $7)",
		user.Email, hashedPassword, user.Username, time.Now(), user.LanguageLevel, user.NativeLanguage, time.Now())
	if err != nil {
		log.Fatal("Error inserting data:", err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Signup successful",
	})
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func getUserPassword(email string) (string, error) {
	conn, err := db.GetDBConnection()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}

	var password string

	err = conn.QueryRow(context.Background(), "SELECT password_hash FROM users WHERE email = $1", email).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("error checking user existence: %w", err)
	}

	return password, nil
}

func passwordComparison(passUser, passDB string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passDB), []byte(passUser))
	return err == nil
}
