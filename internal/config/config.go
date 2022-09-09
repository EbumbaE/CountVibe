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
			"refresh": 		"/refresh",
			"registration": "/registration",
			"diary": 		"/diary",
		},
		Server: server.Config{
			Port: ":8080",
			Paths: map[string]string{
			    "static": 	"../../resource/static",
			    "home": 	"../../resource/static/html/home.html",
			},
		},
		Session: session.Config{
			Paths: map[string]string{
		    	"login": 		 "../../resource/static/html/login.html",
	    		"registration":  "../../resource/static/html/registration.html",
			    
			    "user":   		 "../../resource/templates/user.html",
			    "diary":	 "../../resource/templates/diary/diary.html",
			    "item":			 "../../resource/templates/diary/portion.html",
			},
			FormatsPages: map[string]string{
				"profile": "/%s/",
				"diary":   "/%s/diary", 
			},
			JwtKey: map[string][]byte{
				"access": []byte("birimborom"),
				"refresh": []byte("barumbombom"),
			},
		},
		Certificate: certificate.Config{
			CertPath: "../../internal/certificate/cert.pem",
			KeyPath: "../../internal/certificate/key.pem",
		},
		Database: database.Database{
			Info: "host=localhost port=5432 user=Ebumba password=1234 dbname=abobus sslmode=disable",
		},
	}
}