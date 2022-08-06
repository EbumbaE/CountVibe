package config

import(
	"CountVibe/internal/server"
	"CountVibe/internal/certificate"
)

type Config struct{
	Server server.Config `yaml:"server"`
	Certificate certificate.Config `yaml:"certificate"`
}

func NewConfig() Config{
	return Config{
		server.Config{
			Port: ":8080",
			Pages: server.Pages{
				Begin: "/",
				Home: "/home",
				Login: "/login",
				Auth: "/auth",
				Refresh: "/refresh",
				Diary: "/diary",
				Registration: "/registration",
			},
		},
		certificate.Config{
			Certfile: "../../internal/certificate/cert.pem",
			Keyfile: "../../internal/certificate/key.pem",
		},
	}
}