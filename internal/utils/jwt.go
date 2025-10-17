package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(getSecretKey())

func getSecretKey() string {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		secret = "defaultsecretkey"
	}
	return secret
}

func GenerateJWT(userID, userEmail string) (string, error) {
	payload := jwt.MapClaims{
		"user_id": userID,
		"email":   userEmail,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(jwtSecret)
}
