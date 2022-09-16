package server

import (
	"fmt"
	"net/http"

	"CountVibe/internal/log"
)

type Server struct {
	Port   string
	pages  map[string]string
	paths  map[string]string
	Logger log.Logger
}

func NewServer(c Config, cpages map[string]string, logger log.Logger) *Server {
	return &Server{
		Port:   c.Port,
		pages:  cpages,
		paths:  c.Paths,
		Logger: logger,
	}
}

func (s *Server) Run(Certfile, Keyfile string) {
	s.setupServerHandlers()
	s.Logger.Error(http.ListenAndServeTLS(s.Port, Certfile, Keyfile, nil))
}

func (s *Server) beginHandler(w http.ResponseWriter, r *http.Request) {
	pages := s.pages
	http.Redirect(w, r, pages["home"], http.StatusFound)
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, s.paths["home"])
	case "POST":
		pages := s.pages

		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "parse form ", err)
			return
		}
		button := r.FormValue("button")
		if button == "login" {
			http.Redirect(w, r, pages["login"], http.StatusFound)
		}
		if button == "registration" {
			http.Redirect(w, r, pages["registration"], http.StatusFound)
		}
	}
}

func (s *Server) setupServerHandlers() {
	http.HandleFunc(s.pages["begin"], s.beginHandler)
	http.HandleFunc(s.pages["home"], s.homeHandler)

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(s.paths["static"]))))
}
