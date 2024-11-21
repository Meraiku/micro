package config

import "github.com/joho/godotenv"

type Config struct {
	HTTP     HTTP
	Logstash Logstash
	Service  Service
}

func New() Config {
	return Config{
		HTTP:     NewHTTP(),
		Logstash: NewLogstash(),
		Service:  NewService(),
	}
}

func Load(path string) {
	godotenv.Load(path)
}
