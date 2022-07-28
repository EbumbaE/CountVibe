package config

import(
	"CountVibe/internal/server"
)

type Config struct{
	Server server.Config `yaml:"server"`
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

			Certfile: "../../certificate/cert.pem",
			Keyfile: "../../certificate/key.pem",
		},
	}
}