package goauthpkg

// HashStr is an interface consists of 2 methods for manipulating string and hash
type HashStr interface {
	Generate(s string) (string, error)
	Compare(hash string, s string) error
}
