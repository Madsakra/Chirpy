package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch Chirps", err)
		return
	}

	response := make([]Chirp, 0, len(chirps))

	for _, c := range chirps {
		response = append(response, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			USER_ID:   c.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {

	params := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse ID", err)
		return
	}

	chirp, err2 := cfg.db.GetChirp(r.Context(), chirpID)

	if err2 != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find Chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		USER_ID:   chirp.UserID,
	})
}
