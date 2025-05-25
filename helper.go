package main

import (
	"net/http"

	"example.com/m/v2/internal/auth"
	"github.com/google/uuid"
)

func GetAccessToken(r *http.Request, secretKey string) (uuid.UUID, error) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.Nil, err
	}

	verifiedUser, err := auth.ValidateJWT(bearerToken, secretKey)
	if err != nil {
		return uuid.Nil, err
	}

	return verifiedUser, nil
}
