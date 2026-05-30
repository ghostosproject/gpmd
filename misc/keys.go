package misc

import (
	"crypto/rand"
	"fmt"
)

func GenerateAESKey() ([]byte, error) {
	// AES-256 requires a 32-byte key
	key := make([]byte, 32)

	// Read cryptographically secure random bytes into the key slice
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}
