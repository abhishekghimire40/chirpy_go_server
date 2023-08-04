package main

import (
	"encoding/json"
	"net/http"

	"github.com/abhishekghimire40/chirpy_go_server/database"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CreateUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		userBody := database.User{}
		err := decoder.Decode(&userBody)
		if len(userBody.Email) == 0 && len(userBody.Password) == 0 {
			respondWithError(w, 404, "email or password not provided!")
			return
		}
		if err != nil {
			respondWithError(w, 404, "Provide valid request body")
			return
		}
		hashedPassword, err := hashPassword(userBody.Password)
		if err != nil {
			respondWithError(w, 404, "An error occured while proccessing the information. Please try again later!")
			return
		}
		newUser, err1 := db.CreateUser(userBody.Email, hashedPassword)
		if err1 != nil {
			respondWithError(w, 404, "User with provided email already exists")
			return
		}
		setResponse(w, 201, newUser)
	}
}

func loginUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		newUser := database.User{}
		err := decoder.Decode(&newUser)
		if len(newUser.Email) == 0 && len(newUser.Password) == 0 {
			respondWithError(w, 404, "Valid email or password not provided!")
			return
		}
		if err != nil {
			respondWithError(w, 404, "Provid valid email and password!")
			return
		}
		user, exist := db.GetUser(newUser.Email)
		if !exist {
			respondWithError(w, 404, "There is no user with provided email!")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newUser.Password))
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		publicUser := database.PublicUser{
			Id:    user.Id,
			Email: user.Email,
		}
		setResponse(w, 200, publicUser)

	}
}
