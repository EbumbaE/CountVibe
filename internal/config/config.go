package config

import(
	"CountVibe/internal/server"
	"CountVibe/internal/certificate"
)

type Config struct{
	Server server.Config `yaml:"server"`
	Certificate certificate.Config `yaml:"certificate"`
}

func CreateConfig() Config{
	return Config{
		server.Config{
			Port: ":8080",

			Homepage: "/home",
			Loginpage: "/login",
			Authpage: "/auth",
			Refreshpage: "/refresh",
			Diarypage: "/diary",
			Registrationpage: "/registration",
		},
		certificate.Config{
			Certfile: "../../internal/certificate/cert.pem",
			Keyfile: "../../internal/certificate/key.pem",
		},
	}
}