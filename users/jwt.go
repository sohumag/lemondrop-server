package users

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	jwt.RegisteredClaims
	Email string
}

func GenerateJWT(email string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["Email"] = email
	token.Claims = claims
	return token.SignedString(secretKey)
}

// returns email, err
func ValidateJWT(jwtToken string) (string, error) {

	var userClaim UserClaim

	token, err := jwt.ParseWithClaims(jwtToken, &userClaim, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return userClaim.Email, nil
}

func ParseRequestForJWT(c *fiber.Ctx) (string, error) {
	val := c.Request().Header.Peek("Bearer-Token")
	if string(val) == "" {
		return "", fmt.Errorf("no token in header")
	}

	email, err := ValidateJWT(string(val))
	if err != nil {
		// invalid token
		return "", fmt.Errorf("invalid token")
	}

	return email, nil
}
