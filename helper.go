package main

import (
	"errors"
	"net/http"
	"strings"

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

func GetAPIKey(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("invalid authorization header format")
	}

	apiKey := strings.TrimPrefix(authHeader, prefix)
	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		return "", errors.New("API key missing after prefix")
	}

	return apiKey, nil

}
