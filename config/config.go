package config

type Config struct {
	ServerPort string
}

func GetConfig() *Config {
	return &Config{
		ServerPort: "8080",
	}
}
