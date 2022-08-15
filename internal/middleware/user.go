package middleware

import (
	"fmt"
	"net/http"
)

func deleteCookie(w http.ResponseWriter){
    c := &http.Cookie{
        Name: "token",
        MaxAge: -1,
    }
    http.SetCookie(w, c)
}

func UserExit(w http.ResponseWriter, r *http.Request){
    //moveTokenInBlacklist()
    deleteCookie(w)

    url := r.URL.Path
    http.HandleFunc(url, nil)
}

func (mw Middleware) userHandler(w http.ResponseWriter, r *http.Request){

    ok, _ := LoginVerification(r, mw.jwtKey)
    if !ok{
        url := "/login"
        http.Redirect(w, r, url, http.StatusFound)
        return   
    }

    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, mw.paths.User)
        case "POST":
            err := r.ParseForm()
            if err != nil {
                fmt.Println("parse form ", err)
                return
            }
            button := r.FormValue("button")
            if button == "diary"{
                url := r.URL.Path + "/diary"  
                http.Redirect(w, r, url, http.StatusFound)  
            }
            if button == "exit"{
                UserExit(w, r)
                http.Redirect(w, r, "/", http.StatusFound)  
            }
    }
}

func (mw Middleware) diaryHandler(w http.ResponseWriter, r *http.Request){
	 switch r.Method {
        case "GET":    
            http.ServeFile(w, r, mw.paths.Diary)

   	}
}

func NewHandler(path string){
	
}