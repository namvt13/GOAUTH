package goauthcrypto

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashStr implements goauthpkg.Hash
type HashStr struct{}

const deliminator = "||"

// Generate a salted hash for the input string
func (h *HashStr) Generate(s string) (string, error) {
	salt := uuid.New().String()
	saltedBytes := []byte(salt + s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes)
	return hash + deliminator + salt, nil
}

// Compare string to hash
func (h *HashStr) Compare(hash string, s string) error {
	hashParts := strings.Split(hash, deliminator)
	if len(hashParts) != 2 {
		return errors.New("invalid hash, wrong hash format")
	}

	// Convert exisiting hash and new string (plus salt) to byte
	old := []byte(hashParts[0])
	new := []byte(hashParts[1] + s)

	return bcrypt.CompareHashAndPassword(old, new)
}
