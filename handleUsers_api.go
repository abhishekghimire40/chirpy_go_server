package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/abhishekghimire40/chirpy_go_server/internal/auth"
	"github.com/abhishekghimire40/chirpy_go_server/internal/database"
)

type RequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id            int    `json:"id"`
	Email         string `json:"email"`
	Is_Chirpy_Red bool   `json:"is_chirpy_red"`
	Token         string `json:"token"`
	Refresh_Token string `json:"refresh_token"`
}

type UpgradeUserJson struct {
	Event string `json:"event"`
	Data  struct {
		UserID int `json:"user_id"`
	} `json:"data"`
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
		userResponse := UserResponse{
			Id:    user.Id,
			Email: user.Email,
		}
		access_token, err := auth.GenerateJwtToken("access", userResponse.Id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		refresh_token, err := auth.GenerateJwtToken("refresh", userResponse.Id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		userResponse.Token = access_token
		userResponse.Refresh_Token = refresh_token
		userResponse.Is_Chirpy_Red = user.Is_Chirpy_Red
		setResponse(w, 200, userResponse)

	}
}

func getTokenString(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		return "", errors.New("token not provided")
	}
	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer"))
	return tokenString, nil
}

func updateUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		requestData := RequestBody{}
		err := decoder.Decode(&requestData)
		if err != nil {
			respondWithError(w, 404, "Invalid request body")
			return
		}
		tokenString, err := getTokenString(r)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		token, err := auth.ValidateJwtToken(tokenString)
		if err != nil {
			respondWithError(w, 401, err.Error())
			return
		}
		claims, err := auth.GetTokenClaims(token)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		id, _ := strconv.Atoi(claims.Subject)
		issuer := claims.Issuer
		if issuer != "chirpy-access" {
			respondWithError(w, 401, "Unauthorized")
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

// get api key from header
func getAPIKey(r *http.Request) (string, error) {
	apiKey := r.Header.Get("Authorization")
	if len(apiKey) == 0 {
		return "", errors.New("please provide api key")
	}
	apiKey = strings.TrimSpace(strings.TrimPrefix(apiKey, "ApiKey"))
	return apiKey, nil
}

func upgradeUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := getAPIKey(r)

		if err != nil {
			respondWithError(w, 401, err.Error())
			return
		}
		// apiKey value
		polkaApiKey := os.Getenv("POLKA_API_KEY")
		if apiKey != polkaApiKey {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		decoder := json.NewDecoder(r.Body)
		requestData := UpgradeUserJson{}
		err = decoder.Decode(&requestData)
		if err != nil {
			respondWithError(w, 404, "Bad reqeuest")
			return
		}
		if requestData.Event != "user.upgraded" {
			setResponse(w, 200, nil)
			return
		}
		user_id := requestData.Data.UserID
		err = db.UpgradeUser(user_id)
		if err != nil {
			respondWithError(w, 404, err.Error())
			return
		}
		setResponse(w, 200, nil)
	}
}
