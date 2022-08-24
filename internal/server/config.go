package server

type Config struct{
	Port string `yaml:"port"`
	Pages map[string]string `yaml:"pages"`
	Paths map[string]string `yaml:"paths"`
}