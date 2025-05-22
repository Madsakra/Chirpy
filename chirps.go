package main

import "net/http"

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
