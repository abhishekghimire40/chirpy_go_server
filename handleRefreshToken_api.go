package main

import (
	"net/http"
	"strconv"

	"github.com/abhishekghimire40/chirpy_go_server/internal/auth"
	"github.com/abhishekghimire40/chirpy_go_server/internal/database"
)

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

func refreshToken(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token
		passed_token, err := getTokenString(r)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
		}
		// check if token is valid
		valid_token, err := auth.ValidateJwtToken(passed_token)
		if err != nil {
			respondWithError(w, 401, "Unauthroized")
			return
		}
		// check if it is a refresh token or not
		claims, err := auth.GetTokenClaims(valid_token)
		if err != nil {
			respondWithError(w, 401, "Unauthroized")
			return
		}
		issuer := claims.Issuer
		user_id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			respondWithError(w, 401, "Unauthorized")
		}
		if issuer != "chirpy-refresh" {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		// check if token is already revoked
		if _, isRevoked := db.IsRevoked(passed_token); isRevoked {
			respondWithError(w, 401, "Unauthorized")
			return
		}
		token, err := auth.GenerateJwtToken("access", user_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		setResponse(w, http.StatusOK, RefreshTokenResponse{
			Token: token,
		})
	}
}

func revokeToken(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenString(r)
		if err != nil {
			respondWithError(w, 404, "token not passed in header")
			return
		}
		_, err = db.RevokeToken(token)
		if err != nil {
			respondWithError(w, 200, err.Error())
			return
		}
		setResponse(w, 200, nil)
	}
}
