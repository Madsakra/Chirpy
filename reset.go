package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.DeleteAllUsers(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete all users", err)
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))

}
