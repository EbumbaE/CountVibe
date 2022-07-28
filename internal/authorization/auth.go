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

func getPassword(username string) (password string){
	
    password, err := hashPassword("123")
    if err != nil{
        return ""
    }

    return password
}

func checkUser(username string) bool{
	return true
}

func verifyUserPass(username, password string) bool {
	
	hasUser := checkUser(username)
	if !hasUser {
    	return false
  	}
  	
  	rightPass := getPassword(username)
  	if err := bcrypt.CompareHashAndPassword([]byte(rightPass), []byte(password)); err == nil {
    	return true
  	}

  	return false
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

            ok := verifyUserPass(username, password)
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

            http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
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

func saveUser(username string, hashpassword string){
    
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

            if checkUser(username){
                // todo notisfaction
                return
            }

            hash, err := hashPassword(password)
            if err != nil{

                return
            }
            saveUser(username, hash)

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

            http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
    }
}