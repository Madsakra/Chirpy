package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"example.com/m/v2/internal/auth"
	"example.com/m/v2/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Body    string    `json:"body"`
	USER_ID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	USER_ID   uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleVerification(w http.ResponseWriter, r *http.Request) {

	bearerToken, err1 := auth.GetBearerToken(r.Header)

	if err1 != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retrieve bearer token", err1)
		return
	}

	verifiedUser, err2 := auth.ValidateJWT(bearerToken, cfg.secret_key)

	if err2 != nil {
		println(bearerToken)
		respondWithError(w, http.StatusUnauthorized, "Please Login First before posting", err2)
		return
	}

	// read from params first
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return

	} else {
		badWords := map[string]struct{}{
			"kerfuffle": {},
			"sharbert":  {},
			"fornax":    {},
		}
		filteredWord := wordFilter(params.Body, badWords)

		chirp, err := cfg.db.CreateChirps(r.Context(), database.CreateChirpsParams{
			Body:   filteredWord,
			UserID: verifiedUser,
		})
		if err != nil {
			respondWithError(w, http.StatusNotAcceptable, "Couldn't create Chirps", err)
			return
		}

		respondWithJSON(w, http.StatusCreated, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			USER_ID:   verifiedUser,
		})
	}
}

func wordFilter(word string, badWords map[string]struct{}) string {
	split := strings.Split(word, " ")
	for i, w := range split {
		lowerCased := strings.ToLower(w)
		if _, ok := badWords[lowerCased]; ok {
			split[i] = "****" // replace with asterisks
		}
	}
	finalWord := strings.Join(split, " ")
	return finalWord
}
