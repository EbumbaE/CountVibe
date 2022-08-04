package authorization

import (
	"fmt"
	"net/http"
    "time"
	
    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"

    "CountVibe/internal/database"
)

var jwtKey = []byte("my_secret_key")

type User struct{
    username string `json:"username"`
    password string `json:"password"`
}

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
            http.ServeFile(w, r, "../../static/html/login.html")
        case "POST":

            err := r.ParseForm()
            if err != nil {
                fmt.Fprintf(w, "ParseForm() err: %v", err)
                return
            }

            username := r.FormValue("username")
            password := r.FormValue("password")

            ok, err := verifyUserPass(username, password)
            if err != nil{
                fmt.Fprintf(w, err.Error())
                return
            }
            if !ok{
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

            http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)        //mmmmmmmm huita
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

    w.Write([]byte(fmt.Sprintf("You are authorize %s!", claims.Username)))

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

func hashPassword(password string) (string, error){
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, "../../static/html/registration.html")
        case "POST":

            err := r.ParseForm()
            if err != nil {
                fmt.Fprintf(w, "ParseForm() err: %v", err)
                return
            }

            username := r.FormValue("username")
            password := r.FormValue("password")

            hasUser, err := database.CheckUserInDB(username)
            if err != nil{
                fmt.Fprintf(w, err.Error())
                return
            }
            if !hasUser{
                // todo notisfaction
                return
            }

            hash, err := hashPassword(password)
            if err != nil{

                return
            }
            database.InsertNewUser(username, hash)

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

            http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)            //mmmmmmmmm huita
    }
}