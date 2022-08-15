package middleware

type WayTo struct{
    Static string
    Login string
    Registration string
    User string
    Diary string
}

type User struct{
    username string `json:"username"`
    password string `json:"password"`
}

type FormatsPath struct{
    Profile string
    Diary string
}

type Config struct{
    Paths WayTo
    FormatsPath FormatsPath

    JwtKey []byte 
}