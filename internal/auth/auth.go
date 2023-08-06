package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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
func GenerateJwtToken(expires_in_seconds int, user_id int) (string, error) {
	expiring_time := time.Now()
	if expires_in_seconds > 0 && expires_in_seconds <= 86400 {
		expiring_time = expiring_time.Add(time.Duration(expires_in_seconds) * time.Second)
	} else {
		expiring_time = expiring_time.Add(24 * time.Hour)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiring_time),
		Issuer:    "Chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   fmt.Sprintf("%d", user_id),
	})

	tokenString, err := token.SignedString([]byte(jwt_key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJwtToken(signedToken string) (int, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_key), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, errors.New("claims not accessible")
	}
	subject, _ := strconv.Atoi(claims.Subject)
	return subject, nil
}
