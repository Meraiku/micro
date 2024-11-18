package v1

import (
	"context"

	"github.com/meraiku/micro/user/internal/models"
	"github.com/meraiku/micro/user/pkg/auth_v1"
)

type AuthService interface {
	Login(ctx context.Context, user *models.User) (*models.Tokens, error)
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Authenticate(ctx context.Context, accessToken string) (*models.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*models.Tokens, error)
}

type GRPCAuthServer struct {
	auth_v1.UnimplementedAuthV1Server

	authService AuthService
}

func NewAuthServer(authService AuthService) *GRPCAuthServer {
	return &GRPCAuthServer{
		authService: authService,
	}
}

func (s *GRPCAuthServer) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.Tokens, error) {

	user, err := models.NewUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	tokens, err := s.authService.Login(ctx, user)
	if err != nil {
		return nil, err
	}

	return FromTokens(tokens), nil
}

func (s *GRPCAuthServer) Register(ctx context.Context, req *auth_v1.RegisterRequest) (*auth_v1.RegisterResponse, error) {

	user, err := models.NewUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	user, err = s.authService.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	out := &auth_v1.RegisterResponse{
		Id:       user.ID.String(),
		Username: user.Name,
	}

	return out, nil
}

func (s *GRPCAuthServer) Authenticate(ctx context.Context, req *auth_v1.AuthenticateRequest) (*auth_v1.Tokens, error) {

	accessToken := req.AccessToken

	tokens, err := s.authService.Authenticate(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return FromTokens(tokens), nil
}

func (s *GRPCAuthServer) Refresh(ctx context.Context, req *auth_v1.RefreshRequest) (*auth_v1.Tokens, error) {

	refreshToken := req.RefreshToken

	tokens, err := s.authService.Refresh(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return FromTokens(tokens), nil
}
