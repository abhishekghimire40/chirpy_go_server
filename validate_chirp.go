package main

import (
	"encoding/json"
	"log"
	"net/http"
)

/*
fuction to validate the incoming chirp and return with a response if the chirp
is valid or not
*/
func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirpBody := chirp{}
	err := decoder.Decode(&chirpBody)
	if err != nil || len(chirpBody.Body) == 0 {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(chirpBody.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	respondValidChirp(w)
}

// function to response if any error occurs
func respondWithError(w http.ResponseWriter, code int, errMsg string) {
	type error struct {
		Error string `json:"error"`
	}
	newErr := error{
		Error: errMsg,
	}
	setResponse(w, newErr)
}

// function to respond if the chirp is valid
func respondValidChirp(w http.ResponseWriter) {
	type validChirp struct {
		Valid bool `json:"valid"`
	}
	valid := validChirp{
		Valid: true,
	}
	setResponse(w, valid)

}

// function to set response
func setResponse(w http.ResponseWriter, res interface{}) {
	dat, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
