package server

import (
	"CountVibe/internal/authorization"
	"CountVibe/internal/log"

	"net/http"
)

type Server struct {
	Port string

	homepage string 
	loginpage string 
	authpage string 
	refreshpage string
	diarypage string
	registrationpage string

	Logger log.Logger
}

func CreateServer(c Config, logger log.Logger) *Server{
	return &Server{
		Port: c.Port,

		homepage: c.Homepage,
		loginpage: c.Loginpage,
		authpage: c.Authpage,
		refreshpage: c.Refreshpage,
		diarypage: c.Diarypage,
		registrationpage: c.Registrationpage,

		Logger: logger,
	}
}

func (s *Server) Run(Certfile string, Keyfile string){

	s.setupServerHandlers()
	s.setupAuthHandlers()
	s.Logger.Error(http.ListenAndServeTLS(s.Port, Certfile, Keyfile, nil))

}

func handler(w http.ResponseWriter, r *http.Request){
	http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
}

func homehandler(w http.ResponseWriter, r *http.Request){
    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, "../../static/html/home.html")
    }
}

func (s *Server) setupServerHandlers(){
	http.HandleFunc("/", handler)
	http.HandleFunc(s.homepage, homehandler)
}


func (s *Server) setupAuthHandlers(){

    http.HandleFunc(s.authpage, authorization.AuthHandler)
    http.HandleFunc(s.loginpage, authorization.LoginHandler)
    http.HandleFunc(s.refreshpage, authorization.RefreshHandler)
    http.HandleFunc(s.registrationpage, authorization.RegistrationHandler)
             
}