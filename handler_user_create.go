package main

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/m/v2/internal/auth"
	"example.com/m/v2/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Token         string    `json:"token,omitempty"`
	Refresh_Token string    `json:"refresh_token,omitempty"`
	Is_Chirpy_Red bool      `json:"is_chirpy_red"`
}

type userParams struct {
	Email              string `json:"email"`
	Password           string `json:"password"`
	Expires_In_Seconds int32  `json:"expires_in_seconds"`
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

	passHash, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't Hashpassword", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: passHash,
	})

	if err != nil {
		respondWithError(w, http.StatusNotAcceptable, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:            user.ID,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Email:         user.Email,
		Is_Chirpy_Red: user.IsChirpyRed.Bool,
	})
}
