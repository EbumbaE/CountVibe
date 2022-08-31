package session

import (
    "fmt"
    "net/http"
    "strconv"
    "errors"

    "github.com/dgrijalva/jwt-go"
)

type AuthDetails struct{
    accessUuid string
    userID   int64
}

type User struct{
    id int64 `yaml:"ID"`
    username string `yaml:"username"`
    password string `yaml:"password"`
    isLogin bool `yaml:"islogin"`
}

func deleteCookie(w http.ResponseWriter){
    access := &http.Cookie{
        Name: "access_token",
        MaxAge: -1,
    }
    http.SetCookie(w, access)

    refresh := &http.Cookie{
        Name: "access_token",
        MaxAge: -1,
    }
    http.SetCookie(w, refresh)
}

func (s *Session) getAuthDetails(r *http.Request) (*AuthDetails, error) {
    token, err := GetjwtToken(r, s.jwtKey) 
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if ok && token.Valid {
        
        accessUuid, ok := claims["access_uuid"].(string)
        if !ok {
            return nil, err
        }

        strUserID := fmt.Sprintf("%.f", claims["user_id"])
        UuserID, err := strconv.ParseUint(strUserID, 10, 64)
        if err != nil {
            return nil, err
        }
        userID := int64(UuserID)
        
        ad := &AuthDetails{
            accessUuid: accessUuid,
            userID: userID,
        }
        return ad, nil
    }
    return nil, err
}

func (s *Session)userLogout(w http.ResponseWriter, r *http.Request){
    //moveTokenInBlacklist()
    deleteCookie(w)
    
    ad, err := s.getAuthDetails(r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        s.Logger.Error(err, "get authorization details")
        return
    }
    delErr := s.deleteTokens(ad)
    if delErr != nil {
        w.WriteHeader(http.StatusUnauthorized)
        s.Logger.Error(err, "delete tokens")
        return
    }

    http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Session)verificationUserID(authD *AuthDetails)(bool, error){
    strUserID, err := s.tokensDB.Get(authD.accessUuid).Result()
    if err != nil {
        return false, err
    }
 
    uUserID, _ := strconv.ParseUint(strUserID, 10, 64)
    userID := int64(uUserID)

    if authD.userID != userID {
        return false, errors.New("userID do not mutch")
    }
    return true, nil
}   

func (s *Session)getUser(userID int64)(*User, error){
    user := &User{
        id: 1,
        username: "Ebumba",
        password: "",
        isLogin: true,
    } 
    return user, nil 
}

func (s *Session)setUser(username, password string)(*User, error){
    user := &User{
        id: 1,
        username: "Ebumba",
        password: "",
        isLogin: true,
    } 
    return user, nil 
}


func (s *Session) userHandler(w http.ResponseWriter, r *http.Request){

    ok, err := LoginVerification(r, s.jwtKey)
    if err != nil{
        s.Logger.Error(err, "Login verification")
    }
    if !ok{
        http.Redirect(w, r, s.pages["login"], http.StatusFound)
        return   
    }

    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, s.paths["user"])

        case "POST":
            err := r.ParseForm()
            if err != nil {
                fmt.Println("parse form ", err)
                return
            }
            button := r.FormValue("button")
            if button == "diary"{
                url := r.URL.Path + s.pages["diary"]
                http.Redirect(w, r, url, http.StatusFound)  
            }
            if button == "exit"{
                s.userLogout(w, r)
                http.Redirect(w, r, "/", http.StatusFound)  
            }
    }
}

func (s *Session)diaryHandler(w http.ResponseWriter, r *http.Request){
     switch r.Method {
        case "GET":    
            fmt.Println(s.paths["diary"])
            http.ServeFile(w, r, s.paths["diary"])

    }
}