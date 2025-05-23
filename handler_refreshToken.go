package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"example.com/m/v2/internal/auth"
	"example.com/m/v2/internal/database"
)

type ResponseParams struct {
	Token string `json:"token,omitempty"`
}

func (cfg *apiConfig) CheckRefreshToken(w http.ResponseWriter, r *http.Request) {

	bearerToken, err1 := auth.GetBearerToken(r.Header)

	if err1 != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retrieve bearer token", err1)
		return
	}

	response, err2 := cfg.db.GetUserFromRefreshToken(r.Context(), bearerToken)

	if err2 != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token don't exist", err2)
		return
	}

	now := time.Now()

	if response.ExpiresAt.Before(now) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", errors.New("token expired error"))
		return
	}

	if response.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", errors.New("refresh token revoked"))
		return
	}
	token, err := auth.MakeJWT(response.UserID, cfg.secret_key, 0)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to generate refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ResponseParams{
		Token: token,
	})

}

func (cfg *apiConfig) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	bearerToken, err1 := auth.GetBearerToken(r.Header)

	if err1 != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retrieve bearer token", err1)
		return
	}

	response, err2 := cfg.db.GetUserFromRefreshToken(r.Context(), bearerToken)

	if err2 != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token don't exist", err2)
		return
	}

	now := time.Now()

	if response.ExpiresAt.Before(now) {
		respondWithError(w, http.StatusBadRequest, "Refresh token expired", errors.New("token expired error"))
		return
	}

	if response.RevokedAt.Valid && response.RevokedAt.Time.Before(now) {
		respondWithError(w, http.StatusBadRequest, "Refresh token already revoked", errors.New("refresh token revoked"))
		return
	}

	_, err3 := cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		Token: bearerToken,
	})

	if err3 != nil {
		respondWithError(w, http.StatusInternalServerError, "Failure to revoke token", errors.New("refresh token revoked"))
		return
	}

	w.WriteHeader(204)

}
