package session

type ConfigJwtKey struct {
	Access  string `mapstructure:"access"`
	Refresh string `mapstructure:"refresh"`
}

type Config struct {
	Paths        map[string]string `mapstructure:"paths"`
	FormatsPages map[string]string `mapstructure:"formatsPages"`
	JwtKey       ConfigJwtKey      `mapstructure:"jwtKey"`
}
