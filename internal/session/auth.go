package session

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

func LoginHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, wayToLogin)
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

            token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
            
            tokenString, err := token.SignedString(jwtKey)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            http.SetCookie(w, &http.Cookie{
                Name:    "token",
                Value:   tokenString,
                Expires: expirationTime,
            })

            http.Redirect(w, r, "/auth", http.StatusFound)        //mmmmmmmm huita
    }
}


func hashPassword(password string) (string, error){
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, wayToRegistration)
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

            token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
            
            tokenString, err := token.SignedString(jwtKey)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            http.SetCookie(w, &http.Cookie{
                Name:    "token",
                Value:   tokenString,
                Expires: expirationTime,
            })

            http.Redirect(w, r, "/auth", http.StatusFound)            //mmmmmmmmm huita
    }
}

func AuthHandler(w http.ResponseWriter, r *http.Request){
	
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
            return jwtKey, nil
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

    urlUser := "/" + claims.Username
    http.HandleFunc(urlUser, userHandler)

    urlDiary := urlUser + "/diary"
    http.HandleFunc(urlDiary, diaryHandler)

    http.Redirect(w, r, urlUser, http.StatusFound)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request){

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
            return jwtKey, nil
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
    
    updtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    updtokenString, err := updtoken.SignedString(jwtKey)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   updtokenString,
        Expires: expirationTime,
    })
}