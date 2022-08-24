package authorization

import (
	"fmt"
	"net/http"
    "time"
	
    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"

    "CountVibe/internal/database"
)

type Claims struct{
    Username string `json:"username"`
    jwt.StandardClaims
}

type Authorization struct{
    pages map[string]string
    paths map[string]string
    formatsPages map[string]string
    jwtKey []byte 
}

func NewAuthorization(c Config, confpages map[string]string) *Authorization{
    return &Authorization{
        paths: c.Paths,
        jwtKey: c.JwtKey,
        formatsPages: c.FormatsPages,
        pages: confpages,
    }
}

func (a *Authorization) setupAuthorizationHandlers(){
    pages := a.pages
    http.HandleFunc(pages["auth"], a.authHandler)
    http.HandleFunc(pages["login"], a.loginHandler)
    http.HandleFunc(pages["registration"], a.registrationHandler)
    http.HandleFunc(pages["refresh"], a.refreshHandler)         
}

func (a Authorization) Run(){
    a.setupAuthorizationHandlers()
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

func LoginVerification(r *http.Request, jwtKey []byte) (bool, *Claims){
    c, err := r.Cookie("token")
    if err != nil {
        return false, &Claims{}
    }

    tokenString := c.Value

    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, 
        func(token *jwt.Token) (any, error) {
            return jwtKey, nil
        })

    if err != nil || !token.Valid {
        return false, &Claims{}
    }
    
    return true, claims
}

func createToken(claims *Claims, jwtKey []byte)(string, error){

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (a *Authorization) loginHandler(w http.ResponseWriter, r *http.Request){

    ok, claims := LoginVerification(r, a.jwtKey)
    if ok{
        url := "/" + claims.Username
        http.Redirect(w, r, url, http.StatusFound)
    }

    switch r.Method {
        case "GET": 
            http.ServeFile(w, r, a.paths["login"])
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
                fmt.Fprintf(w, "Verify error" , err.Error())
                return
            }
            if !ok{
                   fmt.Fprintf(w, "Incorrect login or password")
                   w.WriteHeader(http.StatusUnauthorized)
                return
            }

            expirationTime := time.Now().Add(5 * time.Minute)
            claims := &Claims{
                Username: username,
                StandardClaims: jwt.StandardClaims{
                    ExpiresAt: expirationTime.Unix(),
                },
            }

            token, err := createToken(claims, a.jwtKey)
            if err != nil{
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            http.SetCookie(w, &http.Cookie{
                Name:    "token",
                Value:   token,
                Expires: expirationTime,
                HttpOnly: true,
            })

            http.Redirect(w, r, a.pages["auth"], http.StatusFound) 
    }
}

func hashPassword(password string) (string, error){
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func (a *Authorization) registrationHandler(w http.ResponseWriter, r *http.Request){

    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, a.paths["registration"])
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
                fmt.Println("check in db", err.Error())
                return
            }
            if hasUser{
                fmt.Fprintf(w, "User already exist")                   
                // todo notisfaction
                return
            }

            hash, err := hashPassword(password)
            if err != nil{
                fmt.Println("hash password", err.Error())
                return
            }

            if err := database.InsertNewUser(username, hash); err != nil{
                fmt.Println(err.Error())
            }

            expirationTime := time.Now().Add(5 * time.Minute)
            claims := &Claims{
                Username: username,
                StandardClaims: jwt.StandardClaims{
                    ExpiresAt: expirationTime.Unix(),
                },
            }

            token, err := createToken(claims, a.jwtKey)
            if err != nil{
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            http.SetCookie(w, &http.Cookie{
                Name:    "token",
                Value:   token,
                Expires: expirationTime,
                HttpOnly: true,
            })

            http.Redirect(w, r, a.pages["auth"], http.StatusFound)            
    }
}

func (a *Authorization) authHandler(w http.ResponseWriter, r *http.Request){

    c, err := r.Cookie("token")
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    tokenString := c.Value

    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, 
        func(token *jwt.Token) (any, error) {
            return a.jwtKey, nil
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

    formats := a.formatsPages
            
    urlProfile := fmt.Sprintf(formats["profile"], claims.Username)
    http.HandleFunc(urlProfile, a.userHandler)

    urlDiary := fmt.Sprintf(formats["diary"], claims.Username)
    http.HandleFunc(urlDiary, a.diaryHandler)
    
    http.Redirect(w, r, urlProfile, http.StatusFound)

}

func updateToken(claims *Claims, jwtKey []byte) (string, error){
    
    updtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    updtokenString, err := updtoken.SignedString(jwtKey)
    if err != nil {
        return "", err
    }
    return updtokenString, nil
}

func (a *Authorization) refreshHandler(w http.ResponseWriter, r *http.Request){

    c, err := r.Cookie("token")
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    tokenString := c.Value
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, 
        func(token *jwt.Token) (interface{}, error) {
            return a.jwtKey, nil
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
    
    if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    expirationTime := time.Now().Add(5 * time.Minute)
    claims.ExpiresAt = expirationTime.Unix()

    updToken, err := updateToken(claims, a.jwtKey)
    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
    }

    http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   updToken,
        Expires: expirationTime,
        HttpOnly: true,
    })
}