package config

import(
	"CountVibe/internal/server"
	"CountVibe/internal/certificate"
	"CountVibe/internal/database"
)

type Config struct{
	Server server.Config `yaml:"server"`
	Certificate certificate.Config `yaml:"certificate"`
	Database database.Database
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
		database.Database{
			Info: "host=localhost port=5432 user=postgres password=1234 dbname=abobus sslmode=disable",
		},
	}
}