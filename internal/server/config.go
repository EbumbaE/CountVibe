package server

type Config struct {
	Port  string            `mapstructure:"port"`
	Paths map[string]string `mapstructure:"paths"`
}
