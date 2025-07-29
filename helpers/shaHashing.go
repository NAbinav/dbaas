package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword hashes the input password using SHA-256 and returns the hex representation.
func HashPasswordSHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hasher.Sum(nil)
	return hex.EncodeToString(hashedPassword)
}

// CheckPasswordHash compares the input password with the hashed password and returns true if they match.
func CheckPasswordHashSHA256(password, hash string) bool {
	hashedInput := HashPasswordSHA256(password)
	return hashedInput == hash
}
