package config

import (
	"net"
	"os"
)

type cfgREST struct {
	port string
	host string
}

func NewREST() (cfgREST, error) {
	cfg := cfgREST{}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")

	cfg.port = port
	cfg.host = host

	return cfg, nil
}

func (c cfgREST) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
