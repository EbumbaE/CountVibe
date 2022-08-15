package config

import(
	"CountVibe/internal/server"
	"CountVibe/internal/certificate"
	"CountVibe/internal/database"
	"CountVibe/internal/middleware"
)

type Config struct{
	Pages map[string]string `yaml:"pages"`
	Server server.Config `yaml:"server"`
	Middleware middleware.Config `yaml:"middleware"`
	Certificate certificate.Config `yaml:"certificate"`
	Database database.Database
}

func NewConfig() Config{
	return Config{
		map[string]string{
			"begin": "/",
			"home": "/home",
			"login": "/login",
			"auth": "/auth",
			"refresh": "/refresh",
			"registration": "/registration",
		},
		server.Config{
			Port: ":8080",
		},
		middleware.Config{
			Paths: middleware.WayTo{
			    Static: "../../static/",
		    	Login: "../../static/html/login.html",
	    		Registration: "../../static/html/registration.html",
			    User: "../../static/html/user.html",
			    Diary: "../../static/html/diary.html",
			},
			FormatsPath: middleware.FormatsPath{
				Profile: "/%s",
				Diary: "/%s/diary", 
			},
			JwtKey: []byte("my_secret_key"),
		},
		certificate.Config{
			Certfile: "../../internal/certificate/cert.pem",
			Keyfile: "../../internal/certificate/key.pem",
		},
		database.Database{
			Info: "host=localhost port=5432 user=postgres password=1234 dbname=abobus sslmode=disable",
		},
	}
}