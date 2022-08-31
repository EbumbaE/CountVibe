package session

type Config struct{
    Paths map[string]string
    FormatsPages map[string]string
    JwtKey map[string][]byte
}