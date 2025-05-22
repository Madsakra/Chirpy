package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type userParams struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {

	// read from params first
	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusNotAcceptable, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
