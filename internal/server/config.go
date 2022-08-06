package server

type Config struct{
	Port string `yaml:"port"`
	Pages Pages `yaml:"pages"`
}