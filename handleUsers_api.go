package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abhishekghimire40/chirpy_go_server/internal/auth"
	"github.com/abhishekghimire40/chirpy_go_server/internal/database"
)

type RequestBody struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ExpiresAt int    `json:"expires_in_seconds"`
}

func createUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		userBody := RequestBody{}
		err := decoder.Decode(&userBody)

		if err != nil || len(userBody.Email) == 0 || len(userBody.Password) == 0 {
			respondWithError(w, 404, "Provide valid request body")
			return
		}
		newUser, err1 := db.CreateUser(userBody.Email, userBody.Password)
		if err1 != nil {
			respondWithError(w, 404, err.Error())
			return
		}
		setResponse(w, 201, newUser)
	}
}

func loginUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		userBody := RequestBody{}
		err := decoder.Decode(&userBody)
		if err != nil || len(userBody.Email) == 0 || len(userBody.Password) == 0 {
			respondWithError(w, 404, "Provid valid email and password!")
			return
		}
		user, exist := db.GetUser(userBody.Email)
		if !exist {
			respondWithError(w, 404, "There is no user with provided email!")
			return
		}
		err = auth.CheckPasswordHash(userBody.Password, user.Password)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		userResponse := struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
			Token string `json:"token"`
		}{
			Id:    user.Id,
			Email: user.Email,
		}
		token, err := auth.GenerateJwtToken(userBody.ExpiresAt, userResponse.Id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		userResponse.Token = token
		setResponse(w, 200, userResponse)

	}
}

func updateUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		requestData := RequestBody{}
		err := decoder.Decode(&requestData)
		if err != nil {
			respondWithError(w, 404, "Invalid request body")
		}
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer"))
		id, err := auth.ValidateJwtToken(tokenString)
		if err != nil {
			respondWithError(w, 401, err.Error())
			return
		}
		hashedPassword, err := auth.HashPassword(requestData.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		user, err := db.UpdateUser(id, requestData.Email, hashedPassword)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		setResponse(w, 200, user)

	}
}
