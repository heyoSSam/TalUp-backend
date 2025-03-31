package taskWord

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	Key      string `json:"key"`
	Sentence string `json:"sentence"`
}

type Response struct {
	Words []string `json:"words"`
}

func jwtValidator(token string, c echo.Context) (bool, error) {
	fmt.Println("Received token:", token)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return false, errors.New("JWT_SECRET is not set")
	}

	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("JWT parsing error:", err)
		return false, errors.New("invalid or expired token")
	}

	fmt.Println("Token is valid!")
	return true, nil
}

func RegisterRoutes(e *echo.Group) {
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	e.POST("/answers", AnswersGenerator)
}

// AnswersGenerator processes the request and returns words
// @Summary Generate words from a given sentence
// @Description Calls the RoBERTa API and returns a shuffled list of words
// @Tags TaskWord
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body Request true "Input request containing key and sentence"
// @Success 200 {object} Response
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /taskWord/answers [post]
func AnswersGenerator(c echo.Context) error {
	var req Request

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	url := os.Getenv("ROBERTA_URL")
	payload := fmt.Sprintf(`{"text": "%s"}`, req.Sentence)

	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to call RoBERTa API"})
	}
	defer resp.Body.Close()

	var apiResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decode response"})
	}

	result := append([]string{req.Key}, apiResponse.Words...)
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })

	return c.JSON(http.StatusOK, result)
}
