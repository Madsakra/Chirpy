package main

import (
	"errors"
	"net/http"
	"strings"

	"example.com/m/v2/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) DeleteSingleChirp(w http.ResponseWriter, r *http.Request) {

	// Check access token first - uuid
	accessToken, err := GetAccessToken(r, cfg.secret_key)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing Access token or malformed token", err)
		return
	}

	// get the chirp ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	chirpIDStr := pathParts[3] // /api/chirps/{chirpID} → index 3

	if len(chirpIDStr) == 0 {
		println(chirpIDStr)
		respondWithError(w, http.StatusForbidden, "Parameters do not contain chirp ID", errors.New("empty Parameter"))
		return
	}

	// convert id to uuid
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}

	// Step 1: Fetch chirp to see if it exists and who owns it
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {

		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	// Step 2: Check ownership
	if chirp.UserID != accessToken {
		respondWithError(w, http.StatusForbidden, "You are not the author of this chirp", errors.New("unauthorized delete"))
		return
	}

	deleteErr := cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpID,
		UserID: accessToken,
	})

	if deleteErr != nil {
		respondWithError(w, http.StatusNotFound, "Chirp Not Found", deleteErr)
	}

	// success
	w.WriteHeader(204)

}
