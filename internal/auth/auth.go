package auth

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Group) {
	e.POST("/login", Login)
	e.POST("/register", Register)
}

func Login(c echo.Context) error {
	return nil
}

func Register(c echo.Context) error {
	
	return nil
}
