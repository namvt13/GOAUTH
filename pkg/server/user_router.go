package goauthserver

import (
	goauthpkg "chotot/go_auth/pkg"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/sahilm/fuzzy"

	"github.com/gorilla/mux"
)

// userRouter represents the user route and it will holds the user collection for manipulating the user collection
type userRouter struct {
	UserCollection goauthpkg.UserCollection
	UserSession    goauthpkg.UserSession
}

// NewUserRouter used to create new user router
func NewUserRouter(u goauthpkg.UserCollection, router *mux.Router, rc goauthpkg.UserSession) *mux.Router {
	userRouter := userRouter{u, rc}

	router.HandleFunc("/register", userRouter.registerUserHandler).Methods("POST")
	router.HandleFunc("/current", userRouter.getCurrentUserHandler).Methods("GET")
	router.HandleFunc("/login", userRouter.loginHandler).Methods("POST")
	router.HandleFunc("/logout", userRouter.logoutHandler).Methods("GET")
	router.HandleFunc("/refresh", userRouter.refreshHandler).Methods("GET")
	router.HandleFunc("/search/all", userRouter.getAllUserHandler).Methods("GET")
	router.HandleFunc("/search/{username}", userRouter.getUserHandler).Methods("GET")

	return router
}

func decodeUser(r *http.Request) (goauthpkg.UserObj, error) {
	var u goauthpkg.UserObj
	if r.Body == nil {
		return u, errors.New("no request body found")
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	return u, err
}

// createUserHandler will process request coming from PUT "/"
func (ur *userRouter) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUser(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload (invalid JSON format)")
		return
	}

	err = ur.UserCollection.CreateUser(&user)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = ur.UserSession.InsertUser(user.Username)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
	}

	http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)

	// If succeed, send the status OK back to client
}

// getCurrenUserHandler will return current user from database, if token is still not expired
func (ur *userRouter) getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get current token id from cookie
	cookie, err := r.Cookie("go_auth_id")
	if err != nil {
		Error(w, http.StatusUnauthorized, "Cookie is not present")
		return
	}

	sessionToken := cookie.Value
	username, err := ur.UserSession.RetreiveSession(sessionToken)
	if err != nil {
		Error(w, http.StatusUnauthorized, "Session expired")
		return
	}

	log.Println("USERNAME GET: ", username)
	userObj, err := ur.UserCollection.GetByUsername(username)
	if err != nil {
		Error(w, http.StatusInternalServerError, "Can't retreive user from MongoDB")
		return
	}

	// Send back user to client
	JSON(w, http.StatusOK, userObj.Username)
}

// getUserHandler will process request coming from GET "/{username}"
func (ur *userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get all the variable from request (not body, more like params)
	log.Println(r.URL)
	vars := mux.Vars(r)
	log.Println(vars)

	username := vars["username"]
	log.Println("USERNAME GET: ", username)

	// userObj, err := ur.UserCollection.GetByUsername(username)
	// if err != nil {
	// 	Error(w, http.StatusNotFound, err.Error())
	// 	return
	// }

	userList, err := ur.UserSession.RetreiveAllUsers()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
	}
	matches := fuzzy.Find(username, userList)

	usersFound := []string{}
	for _, match := range matches {
		usersFound = append(usersFound, match.Str)
	}

	// Return user back to client
	JSON(w, http.StatusOK, usersFound)
}

// getAllUserHandler handles all request coming from "/user/search/all" path
func (ur userRouter) getAllUserHandler(w http.ResponseWriter, r *http.Request) {
	usernameArr, err := ur.UserCollection.GetAllUser()
	if err != nil {
		Error(w, http.StatusInternalServerError, "Can't get all users")
	}

	JSON(w, http.StatusOK, usernameArr)
}

func decodeCredential(r *http.Request) (goauthpkg.Credential, error) {
	var cred goauthpkg.Credential
	if r.Body == nil {
		return cred, errors.New("no request body found")
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cred)
	return cred, err
}

// loginHandler will process request coming from POST /login
func (ur *userRouter) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AUTHENTICATING...")

	cred, err := decodeCredential(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request (invalid JSON format)")
		return
	}

	userObj, err := ur.UserCollection.Login(cred)
	if err != nil {
		log.Println(err)
		Error(w, http.StatusInternalServerError, "Incorrect password")
		return
	}

	// Create new session token
	token, err := ur.UserSession.SaveSession(cred.Username)
	if err != nil {
		Error(w, http.StatusInternalServerError, "Can't save session into Redis")
	}

	log.Println("Token saved: ", token)

	// Set cookie
	ur.UserSession.SetCookie(w, token)

	JSON(w, http.StatusOK, userObj.Username)
}

// logoutHandler handles all request to "/logout"
func (ur *userRouter) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the cookies
	cookie, err := r.Cookie("go_auth_id")
	if err != nil {
		JSON(w, http.StatusOK, "You're not logged in yet")
		return
	}

	sessionToken := cookie.Value
	http.SetCookie(w, &http.Cookie{
		Name:     "go_auth_id",
		Expires:  time.Now(),
		HttpOnly: true,
		SameSite: 1,
	})

	status, err := ur.UserSession.DeleteSession(sessionToken)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("LOGOUHANDLER Status: ", status)

	JSON(w, http.StatusOK, "You're logged out!")
}

// refreshHandler
func (ur *userRouter) refreshHandler(w http.ResponseWriter, r *http.Request) {
	// Get cookie
	cookie, err := r.Cookie("go_auth_id")
	if err != nil {
		Error(w, http.StatusUnauthorized, "Cookie is deleted")
		return
	}

	sessionToken := cookie.Value
	username, err := ur.UserSession.RetreiveSession(sessionToken)
	if err != nil {
		Error(w, http.StatusBadRequest, "Session expired")
		return
	}

	newSessionToken, err := ur.UserSession.SaveSession(username)
	if err != nil {
		Error(w, http.StatusInternalServerError, "Can't save session")
		return
	}

	// Delete old session
	status, err := ur.UserSession.DeleteSession(sessionToken)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("DELETE SESSION Status: ", status)

	// Set new token to cookie
	ur.UserSession.SetCookie(w, newSessionToken)

	JSON(w, http.StatusOK, map[string]string{
		"oldToken": sessionToken,
		"newToken": newSessionToken,
	})
}
