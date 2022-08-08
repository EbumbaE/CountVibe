package session

var jwtKey = []byte("my_secret_key")

const(
    wayToStatic = "../../static/"

    wayToLogin = wayToStatic + "html/login.html"
    wayToRegistration = wayToStatic + "html/registration.html"
    wayToUser = wayToStatic + "html/user.html"
    wayToDiary = wayToStatic + "html/diary.html"
)

type User struct{
    username string `json:"username"`
    password string `json:"password"`
}

type Session struct{
    
}