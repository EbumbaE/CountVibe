package config

import(
	"CountVibe/internal/server"
	"CountVibe/internal/certificate"
	"CountVibe/internal/database"
	"CountVibe/internal/session"
)

type Config struct{
	Pages map[string]string `yaml:"pages"`
	Server server.Config `yaml:"server"`
	Session session.Config `yaml:"middleware"`
	Certificate certificate.Config `yaml:"certificate"`
	Database database.Database
}

func NewConfig() Config{
	return Config{
		Pages: map[string]string{
			"begin": 		"/",
			"home": 		"/home",
			"login": 		"/login",
			"auth": 		"/auth",
			"refresh": 		"/refresh",
			"registration": "/registration",
		},
		Server: server.Config{
			Port: ":8080",
			Paths: map[string]string{
			    "home": "../../static/html/home.html",
			},
		},
		Session: session.Config{
			Paths: map[string]string{
			    "static": 		"../../static/",
		    	"login": 		"../../static/html/login.html",
	    		"registration": "../../static/html/registration.html",
			    "user": 		"../../static/html/user.html",
			    "diary": 		"../../static/html/diary.html",
			},
			FormatsPages: map[string]string{
				"profile": "/%s",
				"diary": "/%s/diary", 
			},
			JwtKey: map[string][]byte{
				"access": []byte("birimborom"),
				"refresh": []byte("barumbombom"),
			},
		},
		Certificate: certificate.Config{
			Certfile: "../../internal/certificate/cert.pem",
			Keyfile: "../../internal/certificate/key.pem",
		},
		Database: database.Database{
			Info: "host=localhost port=5432 user=postgres password=1234 dbname=abobus sslmode=disable",
		},
	}
}