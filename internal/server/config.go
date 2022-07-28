package server

type Config struct{
	Port string `yaml:"port"`

	Homepage string `yaml:"homepage"`
	Loginpage string `yaml:"loginpage"`
	Authpage string `yaml:"authpage"`
	Refreshpage string `yaml:"refreshpage"`
	Diarypage string `yaml:"diarypage"`
	Registrationpage string `yaml:"registrationpage"`

	Certfile string `yaml:"certfile"`
	Keyfile string `yaml:"keyfile"`
}