package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type UpgradeData struct {
	UserID string `json:"user_id"`
}

type UpgradeParams struct {
	Event string      `json:"event"`
	Data  UpgradeData `json:"data"`
}

func (cfg *apiConfig) UpgradeUser(w http.ResponseWriter, r *http.Request) {

	inputKey, err1 := GetAPIKey(r.Header)

	if err1 != nil {
		println(err1.Error())
		w.WriteHeader(401)
		return
	}

	if inputKey != cfg.polka_key {
		w.WriteHeader(401)
		return
	}

	// read from params first
	decoder := json.NewDecoder(r.Body)
	params := UpgradeParams{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// OTHER EVENTS
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	// convert id to uuid
	userId, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	errUpgrade := cfg.db.UpgradeChirp(r.Context(), userId)

	if errUpgrade != nil {
		respondWithJSON(w, 404, "User can't be found")
	}

	w.WriteHeader(204)

}
