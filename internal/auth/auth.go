package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwt_key = os.Getenv("JWT_SECRET")

// check if the hashed password
func HashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// check if the password matches or not
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// generate jwt token
func GenerateJwtToken(token_type string, user_id int) (string, error) {
	var expiring_time time.Time
	var issuer string
	if token_type == "access" {
		issuer = "chirpy-access"
		expiring_time = time.Now().Add(1 * time.Hour)
	} else if token_type == "refresh" {
		issuer = "chirpy-refresh"
		expiring_time = time.Now().Add(24 * time.Hour)
	} else {
		return "", errors.New("token type not specified")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiring_time),
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   fmt.Sprintf("%d", user_id),
	})

	tokenString, err := token.SignedString([]byte(jwt_key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJwtToken(signedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_key), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token, nil
}

func GetTokenClaims(token *jwt.Token) (*jwt.RegisteredClaims, error) {
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("claims not accessible")
	}
	return claims, nil
}
