package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/abhishekghimire40/chirpy_go_server/internal/auth"
	"github.com/abhishekghimire40/chirpy_go_server/internal/database"
	"github.com/go-chi/chi/v5"
)

type errorMsg struct {
	Error string `json:"error"`
}

// function to process get request to get all the chirps present in the database
func GetAllChirps(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := db.GetChirps()
		if err != nil {
			log.Fatal(err)
		}
		// sorting the chirps according to its id
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })
		setResponse(w, http.StatusOK, chirps)
	}
}

// function to sort our chirps
func sortChirps(chirps []database.Chirp) []database.Chirp {
	sorted := chirps
	sort.Slice(sorted, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })
	return sorted
}

/*
function to validate the incoming chirp from a  Post request, save the chirp
if it  is valid and return with a response of that chirp
*/
func ValidateChirp(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decoding the json request
		decoder := json.NewDecoder(r.Body)
		chirpBody := database.Chirp{}
		err := decoder.Decode(&chirpBody)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}
		// validate the chirp
		if len(chirpBody.Body) > 140 {
			respondWithError(w, 400, "Chirp is too long")
			return
		}
		chirpBody.Body = removeProfanedWords(chirpBody.Body)

		// validate user with jwt
		tokenString, err := getTokenString(r)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		token, err := auth.ValidateJwtToken(tokenString)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}

		claims, err := auth.GetTokenClaims(token)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		user_id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		// create a new chirp with id
		newChirp, err := db.CreateChirp(chirpBody.Body, user_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "new chirp cannot be created")
			return
		}
		setResponse(w, 201, newChirp)
	}
}

// function to remove profaned words in coming response body
func removeProfanedWords(str string) string {
	profanedKeywords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	splittedString := strings.Split(str, " ")
	for i, val := range splittedString {
		_, ok := profanedKeywords[strings.ToLower(val)]
		if !ok {
			continue
		}
		splittedString[i] = "****"
	}
	finalString := strings.Join(splittedString, " ")
	return finalString
}

// request handlerFunc for /api/chirps/{chirpID}
func GetSingleChirp(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get chirpID from url and convert it into int
		id, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
		if err != nil {
			setResponse(w, 404, errorMsg{
				Error: "Provide Valid ChirpID",
			})
			return
		}
		// get all chirps from database
		data, err := db.GetChirps()
		if err != nil {
			setResponse(w, 404, errorMsg{
				Error: "Error",
			})
			return
		}
		// check if chirp is available or not
		if len(data) < id {
			setResponse(w, 404, errorMsg{
				Error: "chirp with provided id not available",
			})
			return
		}
		sortedData := sortChirps(data)
		setResponse(w, 200, sortedData[id-1])
	}
}

// function to handle deleting a chirp
func deleteChirp(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// validate
		tokenString, err := getTokenString(r)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		token, err := auth.ValidateJwtToken(tokenString)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		claims, err := auth.GetTokenClaims(token)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		user_id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// get chirp from database
		chirp_id, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
		if err != nil {
			respondWithError(w, 401, "Provide valid ChirpID")
			return
		}
		statusCode, err := db.DeleteChirp(chirp_id, user_id)
		if err != nil {
			respondWithError(w, statusCode, err.Error())
			return
		}
		setResponse(w, statusCode, nil)
	}
}
