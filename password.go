package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Combine password and salt
	saltedPassword := append([]byte(password), salt...)

	// Create hash
	hash := sha256.Sum256(saltedPassword)

	// Combine salt and hash for storage
	// Format: salt + hash
	result := append(salt, hash[:]...)

	// Return as base64 encoded string
	return base64.StdEncoding.EncodeToString(result), nil
}

func VerifyPassword(hashedPassword, password string) error {
	// Decode the stored password
	decoded, err := base64.StdEncoding.DecodeString(hashedPassword)
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
	saltedPassword := append([]byte(password), salt...)
	newHash := sha256.Sum256(saltedPassword)

	// Compare
	for i := 0; i < len(originalHash); i++ {
		if originalHash[i] != newHash[i] {
			return errors.New("passwords don't match")
		}
	}

	return nil
}