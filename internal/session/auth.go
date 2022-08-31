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
    user *User
    pages map[string]string
    paths map[string]string
    formatsPages map[string]string
    jwtKey map[string][]byte

    Logger log.Logger
    tokensDB *redis.Client
}

func NewSession(c Config, confpages map[string]string, logger log.Logger) *Session{
    return &Session{
        user: nil,
        paths: c.Paths,
        jwtKey: c.JwtKey,
        formatsPages: c.FormatsPages,
        pages: confpages,
        Logger: logger,
    }
}

func (s *Session) setupDefaultHandlers(){
    pages := s.pages
    http.HandleFunc(pages["auth"], s.authHandler)
    http.HandleFunc(pages["login"], s.loginHandler)
    http.HandleFunc(pages["registration"], s.registrationHandler)
    http.HandleFunc(pages["refresh"], s.refreshHandler)     
}

func (s *Session) rebornUserHandlers(){
    
    formats := s.formatsPages
    
    users := []User{                                              //test
        {
            username: "Ebumba",
        },
    }

    for _, user := range users{
        urlProfile := fmt.Sprintf(formats["profile"], user.username)
        http.HandleFunc(urlProfile, s.userHandler)

        urlDiary := fmt.Sprintf(formats["diary"], user.username)
        http.HandleFunc(urlDiary, s.diaryHandler)
    }
}

func (s *Session) setupSessionHandlers(){
    s.setupDefaultHandlers()
    s.rebornUserHandlers()
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
	
	if hasUser, err := database.CheckUserInDB(username); err != nil || !hasUser{
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

func LoginVerification(r *http.Request, jwtKey map[string][]byte) (bool, error){
    c, ok := r.Cookie("access_token")
    if ok != nil {
        return false, nil
    }

    accessToken := c.Value

    token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtKey["access"], nil
    })
    if err != nil{
        return false, err
    }

    _, ok := token.Claims.(jwt.Claims)
    if !ok || !token.Valid{
        return false, nil
    }
    
    return true, nil
}

func (s *Session) loginHandler(w http.ResponseWriter, r *http.Request){

    ok, err := LoginVerification(r, s.jwtKey)
    if err != nil{
        s.Logger.Error(err, "Login verification")
    }
    if ok{
        ad, err := s.getAuthDetails(r)
        
        user, err := s.getUser(ad.userID)
        if err == nil{
            url := "/" + user.username
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
                return
            }
            if !ok{
                   fmt.Fprintf(w, "Incorrect login or password")
                   w.WriteHeader(http.StatusUnauthorized)
                return
            }

            user, err := s.setUser(username, password)
            if err != nil{
                s.Logger.Error(err, "set user in db")
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
            s.user = user

            tokens,  err := newTokens(user.id, s.jwtKey)
            if err != nil{
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            if err := s.newAuthorization(user.id, tokens); err != nil{
                w.WriteHeader(http.StatusForbidden)
                s.Logger.Error(err, "create auth")
                return
            }

            saveTokensInCookie(w, tokens)

            http.Redirect(w, r, s.pages["auth"], http.StatusFound) 
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

            hasUser, err := database.CheckUserInDB(username)
            if err != nil{
                s.Logger.Error(err, "check user in db")
                return
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

            if err := database.InsertNewUser(username, hash); err != nil{
                s.Logger.Error(err, "insert new user in db")
                return
            }

            user, err := s.setUser(username, password)
            if err != nil{
                s.Logger.Error(err, "set user in db")
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
            s.user = user

            tokens,  err := newTokens(user.id, s.jwtKey)
            if err != nil{
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            saveTokensInCookie(w, tokens)

            http.Redirect(w, r, s.pages["auth"], http.StatusFound)            
    }
}

func (s *Session) authHandler(w http.ResponseWriter, r *http.Request){

    c, err := r.Cookie("access_token")
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    accessToken := c.Value

    token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.jwtKey["access"], nil
    })
    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if !token.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }   

    formats := s.formatsPages
    urlProfile := fmt.Sprintf(formats["profile"], s.user.username)
    
    http.Redirect(w, r, urlProfile, http.StatusFound)

}

func (s *Session) refreshHandler(w http.ResponseWriter, r *http.Request){

    c, err := r.Cookie("refresh_token")
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    refreshToken := c.Value

    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error){
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.jwtKey["refresh"], nil
    })
    if err != nil{
        s.Logger.Error(err, "Parse token")
        w.WriteHeader(http.StatusUnauthorized)
        return
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
    uUserID, convErr := strconv.ParseUint(strUserID, 10, 64)
    if convErr != nil {
        w.WriteHeader(http.StatusUnprocessableEntity)
        s.Logger.Error(convErr, "convertation userID")
        return
    }
    userID := int64(uUserID)

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
    
    saveTokensInCookie(w, tokens)
    
}