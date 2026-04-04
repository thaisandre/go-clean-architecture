package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	return &WebServer{
		Router:        router,
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	s.Handlers[method+path] = handler
	switch method {
	case "GET":
		s.Router.Get(path, handler)
	case "POST":
		s.Router.Post(path, handler)
	}
}

// start the server
func (s *WebServer) Start() {
	http.ListenAndServe(s.WebServerPort, s.Router)
}
