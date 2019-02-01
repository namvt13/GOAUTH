package goauthmock

// HashStr is a mocking empty struct
type HashStr struct{}

// Generate always return true
func (h *HashStr) Generate(s string) (string, error) {
	return s, nil
}

// Compare always return true
func (h *HashStr) Compare(hash string, s string) error {
	return nil
}
