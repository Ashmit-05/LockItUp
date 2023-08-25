package middlewares

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func GenerateToken(userId string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in reading .env file")
	}
	sK := os.Getenv("SECRET_KEY")
	var secretKey = []byte(sK)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(7*24 * time.Hour).Unix()
	claims["authorized"] = true
	claims["userId"] = userId
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
