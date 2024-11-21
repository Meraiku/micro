package config

import (
	"net"
	"os"
)

type HTTP struct {
	Host string
	Port string
}

func NewHTTP() HTTP {
	host := os.Getenv("HOST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "2000"
	}
	return HTTP{
		Host: host,
		Port: port,
	}
}

func (h HTTP) Addr() string {
	return net.JoinHostPort(h.Host, h.Port)
}
