package main

import (
	"encoding/json"
	"net/http"

	"example.com/m/v2/internal/auth"
)

func (cfg *apiConfig) Login(w http.ResponseWriter, r *http.Request) {

	// read from params first
	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	userDetails, err := cfg.db.GetUserHash(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed retreive password hash", err)
		return
	}

	err2 := auth.CheckPasswordHash(userDetails.HashedPassword, params.Password)
	if err2 != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong password or username", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        userDetails.ID,
		CreatedAt: userDetails.CreatedAt,
		UpdatedAt: userDetails.UpdatedAt,
		Email:     userDetails.Email,
	})

}
