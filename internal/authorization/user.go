package authorization

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


func (a Authorization)exitHandler(w http.ResponseWriter, r *http.Request){
    http.Redirect(w, r, a.pages["home"], http.StatusFound)
}

func (a Authorization)userExit(w http.ResponseWriter, r *http.Request){
    //moveTokenInBlacklist()
    deleteCookie(w)
    
    //url := r.URL.Path
    //http.HandleFunc(url, a.exitHandler)
}

func (a Authorization) userHandler(w http.ResponseWriter, r *http.Request){

    ok, _ := LoginVerification(r, a.jwtKey)
    if !ok{
        http.Redirect(w, r, a.pages["login"], http.StatusFound)
        return   
    }

    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, a.paths["user"])

        case "POST":
            err := r.ParseForm()
            if err != nil {
                fmt.Println("parse form ", err)
                return
            }
            button := r.FormValue("button")
            if button == "diary"{
                url := r.URL.Path + a.pages["diary"]
                http.Redirect(w, r, url, http.StatusFound)  
            }
            if button == "exit"{
                a.userExit(w, r)
                http.Redirect(w, r, "/", http.StatusFound)  
            }
    }
}

func (a Authorization) diaryHandler(w http.ResponseWriter, r *http.Request){
     switch r.Method {
        case "GET":    
            fmt.Println(a.paths["diary"])
            http.ServeFile(w, r, a.paths["diary"])

    }
}