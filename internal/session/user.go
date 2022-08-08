package session

import (
	"fmt"
	"net/http"
)

func userHandler(w http.ResponseWriter, r *http.Request){
    switch r.Method {
        case "GET":    
            http.ServeFile(w, r, wayToUser)
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
            if button == "deleteuser"{
                //http.Redirect(w, r, "/", http.StatusFound)  
            }
    }
}

func diaryHandler(w http.ResponseWriter, r *http.Request){
	 switch r.Method {
        case "GET":    
            http.ServeFile(w, r, wayToDiary)

   	}
}

func NewHandler(path string){
	
}