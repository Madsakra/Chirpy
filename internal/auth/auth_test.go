package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	// Generate a test user ID
	userID := uuid.New()
	tokenSecret := "mySecretKey"
	expiresIn := 30 * time.Minute

	// Create JWT
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	// Validate JWT
	returnedID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if returnedID != userID {
		t.Errorf("Expected userID %v, got %v", userID, returnedID)
	}
}
