package goauthserver

import (
	goauthpkg "chotot/go_auth/pkg"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

// Server represents a server object with router
type Server struct {
	router *mux.Router
}

// NewServer create a new server with user collection service
func NewServer(u goauthpkg.UserCollection, rc goauthpkg.UserSession) *Server {
	s := Server{
		router: mux.NewRouter(),
	}

	// "/user" will be the main path, any subsequent path registered on this subrouter will add its path to "/user"
	NewUserRouter(u, s.newSubrouter("/user"), rc)
	return &s
}

// Start activate new server and listen to port 4040
func (s *Server) Start() {
	log.Println("Server is listening @ http://localhost:4040")
	err := http.ListenAndServe(":4040", handlers.LoggingHandler(os.Stdout, s.router))
	if err != nil {
		log.Fatal("Error while starting server: ", err)
	}
}

// newSubrouter used to register new subrouter
func (s *Server) newSubrouter(path string) *mux.Router {
	return s.router.PathPrefix(path).Subrouter()
}
