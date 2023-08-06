package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
