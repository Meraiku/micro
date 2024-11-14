package config

import (
	"net"
	"os"
)

type GRPC struct {
	port         string
	host         string
	userRepoType string
}

func NewGRPC() (*GRPC, error) {
	cfg := &GRPC{}

	port := os.Getenv("PORT")
	if port == "" {
		port = "20001"
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

func (c *GRPC) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

func (c *GRPC) UserRepoType() string {
	return c.userRepoType
}
