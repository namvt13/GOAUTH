package goauthcrypto_test

import (
	goauthcrypto "chotot/go_auth/pkg/crypto"
	"testing"
)

func Test_Hash(t *testing.T) {
	t.Run("Can hash and compare", hashAndCompare)
	t.Run("Can detect unequal hashes", unequalHashCompare)
	t.Run("Generates a unique salt everytime", generateSalt)
}

func hashAndCompare(t *testing.T) {
	// Arrange
	h := goauthcrypto.HashStr{}
	testString := "testString"

	// Act & Assert
	hash, hashErr := h.Generate(testString)
	if hashErr != nil {
		t.Errorf("Error while generating hash: %s", hashErr)
	}

	compareErr := h.Compare(hash, testString)
	if compareErr != nil {
		t.Errorf("Error while comparing hash and original string: %s", compareErr)
	}

	if hash == testString {
		t.Error("Hash and input string are the same!")
	}
}

func unequalHashCompare(t *testing.T) {
	// Arrange
	h := goauthcrypto.HashStr{}
	testString111 := "testString111"
	testString222 := "testString222"

	// Act & Assert
	hash, hashErr := h.Generate(testString111)
	if hashErr != nil {
		t.Errorf("Error while generating hash: %s", hashErr)
	}

	compareErr := h.Compare(hash, testString222)
	if compareErr == nil {
		t.Errorf("Compare 2 different string successful?")
	}

	if hash == testString111 {
		t.Error("Hash and input string are the same")
	}
}

func generateSalt(t *testing.T) {
	// Arrange
	h := goauthcrypto.HashStr{}
	inputString := "inputString"

	// Act & Assert
	hash1, hashErr := h.Generate(inputString)
	hash2, hashErr := h.Generate(inputString)
	if hashErr != nil {
		t.Errorf("Error while generating hash: %s", hashErr)
	}

	if hash1 == hash2 {
		t.Error("Subsequent hashes must not be the same")
	}
}
