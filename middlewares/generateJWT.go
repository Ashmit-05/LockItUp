package generateJWT

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
	secretKey := os.Getenv("SECRET_KEY")
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * time.Hour)
	claims["authorized"] = true
	claims["userId"] = userId
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
