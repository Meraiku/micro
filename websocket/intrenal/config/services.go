package config

import "os"

type Service struct {
	Auth AuthService
}

type AuthService struct {
	Addr string
}

func NewService() Service {
	return Service{
		Auth: NewAuthService(),
	}
}

func NewAuthService() AuthService {
	addr := os.Getenv("AUTH_SERVICE")
	return AuthService{
		Addr: addr,
	}
}
