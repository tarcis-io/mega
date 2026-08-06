package config

type Config struct {
	ServerHost string
	ServerPort string
}

func Load() *Config {
	return &Config{
		ServerHost: "0.0.0.0",
		ServerPort: "8080",
	}
}
