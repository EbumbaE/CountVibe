package server

type Pages struct {
	Begin string
	Home string 
	Login string 
	Auth string 
	Refresh string
	Diary string
	Registration string
}

type Config struct{
	Port string `yaml:"port"`
	Pages Pages `yaml:"pages"`
}