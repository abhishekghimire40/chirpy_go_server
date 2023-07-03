package database

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type errorMsg struct {
	Error string `json:"error"`
}

// function to process get request to get all the chirps present in the database
func GetAllChirps(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db.mux.RLock()
		chirps, err := db.GetChirps()
		db.mux.RUnlock()
		if err != nil {
			log.Fatal(err)
		}
		setResponse(w, http.StatusOK, chirps)
	}
}

/*
function to validate the incoming chirp from a  Post request, save the chirp
if it  is valid and return with a response of that chirp
*/
func ValidateChirp(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decoding the json request
		decoder := json.NewDecoder(r.Body)
		chirpBody := Chirp{}
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
		// create a new chirp with id
		newChirp, err := db.CreateChirp(chirpBody.Body)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "new chirp cannot be created")
			return
		}
		// ensure that the database is present
		db.mux.RLock()
		db.ensureDB()
		// load data from our database to add the new chirp to our db
		allChirps, err := db.loadDB()
		db.mux.RUnlock()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "previous chirps cannot be laoded!")
			return
		}
		// add the new chirp to the database
		allChirps.Chirps[newChirp.Id] = newChirp
		// save updated chirps to our database
		db.mux.Lock()
		err = db.writeDB(allChirps)
		db.mux.Unlock()
		if err != nil {
			fmt.Println("Error while writing file: ", err)
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

// function to response if any error occurs
func respondWithError(w http.ResponseWriter, code int, errMsg string) {

	newErr := errorMsg{
		Error: errMsg,
	}
	setResponse(w, code, newErr)
}

// function to set response
func setResponse(w http.ResponseWriter, code int, res interface{}) {
	dat, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
