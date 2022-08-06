package server

import (
	"CountVibe/internal/authorization"
	"CountVibe/internal/log"

	"net/http"
)

type Pages struct {
	Begin string
	Home string 
	Login string 
	Auth string 
	Refresh string
	Diary string
	Registration string
}

type Server struct {
	Port string
	Pages Pages
	Logger log.Logger
}

func NewServer(c Config, logger log.Logger) *Server{
	return &Server{
		Port: c.Port,
		Pages: c.Pages,
		Logger: logger,
	}
}

func (s *Server) Run(Certfile string, Keyfile string){

	s.setupServerHandlers()
	s.setupAuthHandlers()
	s.Logger.Error(http.ListenAndServeTLS(s.Port, Certfile, Keyfile, nil))
}

func (s *Server) Shutdown(){

}

func (s *Server) beginHandler(w http.ResponseWriter, r *http.Request){
	pages := s.Pages
	http.Redirect(w, r, pages.Home, http.StatusTemporaryRedirect)
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method {
        case "GET":    
        	way := "../../static/html/home.html"
            http.ServeFile(w, r, way)
    }
}

func (s *Server) setupServerHandlers(){
	pages := s.Pages
	http.HandleFunc(pages.Begin, s.beginHandler)
	http.HandleFunc(pages.Home, s.homeHandler)
}


func (s *Server) setupAuthHandlers(){
	pages := s.Pages
    http.HandleFunc(pages.Auth, authorization.AuthHandler)
    http.HandleFunc(pages.Login, authorization.LoginHandler)
    http.HandleFunc(pages.Refresh, authorization.RefreshHandler)
    http.HandleFunc(pages.Registration, authorization.RegistrationHandler)
             
}