package config

import (
	"net"
	"os"
)

type cfgGRPC struct {
	port string
	host string
}

func NewGRPC() (cfgGRPC, error) {
	cfg := cfgGRPC{}

	port := os.Getenv("PORT")
	if port == "" {
		port = "20001"
	}

	host := os.Getenv("HOST")

	cfg.port = port
	cfg.host = host

	return cfg, nil
}

func (c cfgGRPC) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
