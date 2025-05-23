package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {

	key := make([]byte, 32)

	// always return with no error
	// no need to use the hexEnc val, just use key
	rand.Read(key)

	return hex.EncodeToString(key), nil

}
