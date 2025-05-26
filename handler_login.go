package main

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/m/v2/internal/auth"
	"example.com/m/v2/internal/database"
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

	// GENERATE ACCESS TOKEN
	token, err := auth.MakeJWT(userDetails.ID, cfg.secret_key, time.Duration(params.Expires_In_Seconds))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to generate token", err)
		return
	}

	// GENERATE REFRESH TOKEN
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to generate refresh token", err)
		return
	}

	// store refresh token in db
	_, errRef := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userDetails.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	if errRef != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to store refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:            userDetails.ID,
		CreatedAt:     userDetails.CreatedAt,
		UpdatedAt:     userDetails.UpdatedAt,
		Email:         userDetails.Email,
		Token:         token,
		Refresh_Token: refreshToken,
		Is_Chirpy_Red: userDetails.IsChirpyRed.Bool,
	})

}
