package main

import (
	"encoding/json"
	"net/http"

	"github.com/abhishekghimire40/chirpy_go_server/database"
)

func CreateUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		userBody := database.User{}
		err := decoder.Decode(&userBody)
		if len(userBody.Email) == 0 {
			respondWithError(w, 404, "Valid email not provided!")
			return
		}
		if err != nil {
			respondWithError(w, 404, "Provide valid request body")
			return
		}
		newUser, err1 := db.CreateUser(userBody.Email)
		if err1 != nil {
			respondWithError(w, 404, "User not created! Something went wrong")
			return
		}
		setResponse(w, 201, newUser)
	}
}
