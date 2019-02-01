package goauthpkg

import "net/http"

// UserObj An object represent user in database
type UserObj struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserCollection holds function to manipulate the user collection
type UserCollection interface {
	// Create user in the database
	CreateUser(u *UserObj) error
	GetByUsername(username string) (UserObj, error)
	GetAllUser() (*[]UserObj, error)
	Login(c Credential) (UserObj, error)
}

// UserSession holds methods to use for session authentication
type UserSession interface {
	SaveSession(username string) (string, error)
	RetreiveSession(token string) (string, error)
	DeleteSession(token string) (string, error)
	SetCookie(w http.ResponseWriter, token string)
	InsertUser(username string) error
	RetreiveAllUsers() ([]string, error)
}
