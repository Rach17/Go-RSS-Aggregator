package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func Hash(text string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Combine text and salt
	saltedText := append([]byte(text), salt...)

	// Create hash
	hash := sha256.Sum256(saltedText)

	// Combine salt and hash for storage
	// Format: salt + hash
	result := append(salt, hash[:]...)

	// Return as base64 encoded string
	return base64.StdEncoding.EncodeToString(result), nil
}

func VerifyHash(hashedText, text string) error {
	// Decode the stored text
	decoded, err := base64.StdEncoding.DecodeString(hashedText)
	if err != nil {
		return err
	}

	if len(decoded) < 16+sha256.Size {
		return errors.New("invalid hashed password format")
	}

	// Extract salt and original hash
	salt := decoded[:16]
	originalHash := decoded[16 : 16+sha256.Size]

	// Compute new hash
	saltedText := append([]byte(text), salt...)
	newHash := sha256.Sum256(saltedText)

	// Compare
	for i := 0; i < len(originalHash); i++ {
		if originalHash[i] != newHash[i] {
			return errors.New("texts don't match")
		}
	}

	return nil
}