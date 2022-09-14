package session

import (
	"fmt"
	"net/http"
    "strconv"

    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
    "github.com/go-redis/redis/v7"
    
    "CountVibe/internal/log"
    "CountVibe/internal/database"
)

type Session struct{
    pages map[string]string
    paths map[string]string
    formatsPages map[string]string
    jwtKey map[string][]byte

    idGenerator IDGenerator
    Logger log.Logger
    tokensDB *redis.Client
}

func NewSession(c Config, confpages map[string]string, logger log.Logger) *Session{
    return &Session{
        pages: confpages,
        paths: c.Paths,
        formatsPages: c.FormatsPages,
        jwtKey: c.JwtKey,

        idGenerator: IDGenerator{id: 0},
        Logger: logger,
    }
}

func (s *Session) setupDefaultHandlers(){
    http.HandleFunc(s.pages["login"], s.loginHandler)
    http.HandleFunc(s.pages["registration"], s.registrationHandler)
    http.HandleFunc(s.pages["refresh"], s.refreshHandler)     
}

func (s *Session) restoreUserHandlers(){
    
    usernames, err := database.GetAllUsernames()
    if err != nil{
        s.Logger.Error(err, "get all usernames")
        return
    }

    for _, username := range usernames{
        urlProfile := fmt.Sprintf(s.formatsPages["profile"], username)
        urlDiary := fmt.Sprintf(s.formatsPages["diary"], username)
        
        http.HandleFunc(urlProfile, s.userHandler)
        http.HandleFunc(urlDiary, s.diaryHandler)
    }
}

func (s *Session) setupSessionHandlers(){
    s.setupDefaultHandlers()
    s.restoreUserHandlers()
}

func (s *Session) Run(){
    s.newTokensDB()
    if err := s.checkHealthTokensDB(); err != nil{
        s.Logger.Error("check health TokensDB")
        return
    }
    s.setupSessionHandlers()
}

func verifyUserPass(username, password string)(bool, error) {
	
	if hasUser, err := database.CheckUsernameInDB(username); err != nil || !hasUser{
        return false, err
    }

  	rightPass, err := database.GetUserPassword(username)
    if err != nil{
        return false, err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(rightPass), []byte(password)); err == nil {
    	return true, nil
  	}

  	return false, nil
}

func (s *Session) loginVerification(r *http.Request, jwtKey map[string][]byte) (bool, error){

    token, err := GetjwtToken(r, jwtKey["access"], "access_token")
    if err != nil{
        return false, err
    }

    ok, err := s.confirmToken(token, "access_uuid")

    return ok, err
}

func (s *Session) loginHandler(w http.ResponseWriter, r *http.Request){

    ok, err := s.loginVerification(r, s.jwtKey)
    if err != nil{
        s.Logger.Error(err, "Login verification")
    }
    if ok{
        ad, err := s.getAuthDetails(r)
        
        strUserID := strconv.FormatInt(ad.userID, 10)
        username, err := database.GetUsername(strUserID)
        if err == nil{
            url := "/" + username
            http.Redirect(w, r, url, http.StatusFound)
            return
        }
    }

    switch r.Method {
        case "GET": 
            http.ServeFile(w, r, s.paths["login"])
        case "POST":

            err := r.ParseForm()
            if err != nil {
                fmt.Fprintf(w, "parse form err: %v", err)
                return
            }

            username := r.FormValue("username")
            password := r.FormValue("password")

            ok, err := verifyUserPass(username, password)
            if err != nil{
                s.Logger.Error(err, "verify user password")
            }
            if !ok{
                   fmt.Fprintf(w, "Incorrect login or password")
                   w.WriteHeader(http.StatusUnauthorized)
                return
            }

            strUserID, err := database.GetUserID(username)
            if err !=    nil{
                s.Logger.Error(err, "get userID")
                w.WriteHeader(http.StatusUnauthorized)
                return   
            }
            userID, err := strconv.ParseInt(strUserID, 10, 64)
            if err != nil{
                s.Logger.Error(err, "convertation strUserID to userID")
                w.WriteHeader(http.StatusInternalServerError)
                return   
                
            }

            tokens, err := newTokens(userID, s.jwtKey)
            if err != nil{
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            if err := s.newAuthorization(userID, tokens); err != nil{
                w.WriteHeader(http.StatusForbidden)
                s.Logger.Error(err, "create auth")
                return
            }

            saveTokensCookie(w, tokens)

            urlProfile := fmt.Sprintf(s.formatsPages["profile"], username)
            http.Redirect(w, r, urlProfile, http.StatusFound) 
    }
}

func hashPassword(password string) (string, error){
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func (s *Session) registrationHandler(w http.ResponseWriter, r *http.Request){

    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, s.paths["registration"])
        case "POST":

            err := r.ParseForm()
            if err != nil {
                fmt.Fprintf(w, "parse form err: %v", err)
                return
            }

            username := r.FormValue("username")
            password := r.FormValue("password")

            hasUser, err := database.CheckUsernameInDB(username)
            if err != nil{
                s.Logger.Error(err, "check user in db")
            }
            if hasUser{
                fmt.Fprintf(w, "User already exist")                   
                // todo notisfaction
                return
            }

            hash, err := hashPassword(password)
            if err != nil{
                s.Logger.Error(err, "hash password")
                return
            }

            userID := s.idGenerator.newID()
            strUserID := strconv.FormatInt(userID, 10)

            if err := database.InsertNewUser(strUserID, username, hash); err != nil{
                s.Logger.Error(err, "insert new user in db")
                return
            }

            tokens,  err := newTokens(userID, s.jwtKey)
            if err != nil{
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            if err := s.newAuthorization(userID, tokens); err != nil{
                w.WriteHeader(http.StatusForbidden)
                s.Logger.Error(err, "create auth")
                return
            }

            saveTokensCookie(w, tokens)

            urlProfile := fmt.Sprintf(s.formatsPages["profile"], username)
            urlDiary := fmt.Sprintf(s.formatsPages["diary"], username)
            http.HandleFunc(urlProfile, s.userHandler)
            http.HandleFunc(urlDiary, s.diaryHandler)

            http.Redirect(w, r, urlProfile, http.StatusFound)            
    }
}

func (s *Session) refreshHandler(w http.ResponseWriter, r *http.Request){

    token, err := GetjwtToken(r, s.jwtKey["refresh"], "refresh_token")
    if err != nil{
        s.Logger.Error(err, "get jwt token")
    }

    ok, err := s.confirmToken(token, "refresh_uuid")
    if !ok{
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    if err != nil{
        s.Logger.Error(err, "Confirm token while refresh")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid{
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    
    refreshUuid, ok := claims["refresh_uuid"].(string)
    if !ok{
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    strUserID := fmt.Sprintf("%.f", claims["user_id"])
    userID, convErr := strconv.ParseInt(strUserID, 10, 64)
    if convErr != nil {
        w.WriteHeader(http.StatusUnprocessableEntity)
        s.Logger.Error(convErr, "convertation userID")
        return
    }

    deleted, delErr := s.deleteAuthorization(refreshUuid)
    if delErr != nil || deleted == 0 {
        w.WriteHeader(http.StatusUnauthorized)
        s.Logger.Error(delErr, "delete auth")
        return
    }

    tokens, createErr := newTokens(userID, s.jwtKey)
    if  createErr != nil{
        w.WriteHeader(http.StatusForbidden)
        s.Logger.Error(createErr, "new tokens")
        return
    }

    if err := s.newAuthorization(userID, tokens); err != nil{
        w.WriteHeader(http.StatusForbidden)
        s.Logger.Error(err, "create auth")
        return
    }
    
    saveTokensCookie(w, tokens)
    
}