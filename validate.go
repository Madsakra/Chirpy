package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type cleanedParams struct {
	Cleaned_Body string `json:"cleaned_body"`
}

func handleVerification(w http.ResponseWriter, r *http.Request) {

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

		respondWithJSON(w, http.StatusOK, cleanedParams{
			Cleaned_Body: filteredWord,
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
