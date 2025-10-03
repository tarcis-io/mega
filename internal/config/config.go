package config

const (
	serverAddressEnvKey     = "SERVER_ADDRESS"
	serverAddressEnvDefault = "localhost:8080"
)

type (
	Config struct {
		ServerAddress string
	}
)
