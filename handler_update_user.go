package main

import (
	"encoding/json"
	"net/http"

	"example.com/m/v2/internal/auth"
	"example.com/m/v2/internal/database"
)

type updateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateResponse struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) UpdateUserAccount(w http.ResponseWriter, r *http.Request) {

	// get the access token first
	accessToken, err := GetAccessToken(r, cfg.secret_key)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing Access token or malformed token", err)
		return
	}

	// read from params
	decoder := json.NewDecoder(r.Body)
	params := updateUserParams{}
	err2 := decoder.Decode(&params)

	if err2 != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// hash the password
	newPass, err3 := auth.HashPassword(params.Password)
	if err3 != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	errUpdate := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: newPass,
		ID:             accessToken,
	})

	if errUpdate != nil {
		respondWithError(w, http.StatusInternalServerError, "Update account failed", err)
		return
	}

	respondWithJSON(w, http.StatusOK, updateResponse{
		Email: params.Email,
	})

}
