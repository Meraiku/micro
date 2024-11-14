package config

import (
	"net"
	"os"
)

type REST struct {
	port         string
	host         string
	userRepoType string
}

func NewREST() (*REST, error) {
	cfg := &REST{}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")

	userRepo := os.Getenv("USER_REPO")
	if userRepo == "" {
		userRepo = "memory"
	}

	cfg.port = port
	cfg.host = host
	cfg.userRepoType = userRepo

	return cfg, nil
}

func (c *REST) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

func (c *REST) UserRepoType() string {
	return c.userRepoType
}
