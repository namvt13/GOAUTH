package goauthpkg

// Credential represent the object holding our authentication data
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
