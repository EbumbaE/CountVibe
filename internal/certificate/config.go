package certificate

type Config struct{
	Certfile string `yaml:"certfile"`
	Keyfile string `yaml:"keyfile"`
}