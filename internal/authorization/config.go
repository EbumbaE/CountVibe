package authorization

type User struct{
    username string `json:"username"`
    password string `json:"password"`
    isLogin bool `json:"islogin"`
}

type TokenDetails struct {
    AccessToken string
    AccessUuid string
    AccessExpires int64

    RefreshToken string
    RefreshUuid string
    RefreshExpires int64
}

type Config struct{
    Paths map[string]string
    FormatsPages map[string]string

    JwtKey []byte 
}