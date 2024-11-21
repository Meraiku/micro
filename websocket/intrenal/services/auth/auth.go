package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/pkg/auth_v1"
	"github.com/meraiku/micro/websocket/intrenal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	client auth_v1.AuthV1Client
}

func New(ctx context.Context, addr string) (*Service, error) {

	logging.L(ctx).Info(
		"creating auth service",
		logging.String("addr", addr),
	)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.L(ctx).Error(
			"failed to create grpc client",
			logging.Err(err),
			logging.String("addr", addr),
		)

		return nil, err
	}

	client := auth_v1.NewAuthV1Client(conn)

	logging.L(ctx).Info(
		"auth service created",
	)

	return &Service{
		client: client,
	}, nil
}

func (s *Service) Login(ctx context.Context, user *models.User) (*models.Tokens, error) {

	logging.WithAttrs(ctx,
		logging.String("username", user.Name),
		logging.String("operation", "login"),
	)

	log := logging.L(ctx)

	req := &auth_v1.LoginRequest{
		Username: user.Name,
		Password: user.Password,
	}

	log.Debug(
		"login handler called",
	)

	tokens, err := s.client.Login(ctx, req)
	if err != nil {
		log.Error(
			"failed to login",
			logging.Err(err),
		)

		return nil, err
	}

	log.Debug(
		"login success",
		logging.String("access_token", tokens.AccessToken),
		logging.String("refresh_token", tokens.RefreshToken),
	)

	out := &models.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	return out, nil
}

func (s *Service) Register(ctx context.Context, user *models.User) (*models.User, error) {

	logging.WithAttrs(ctx,
		logging.String("username", user.Name),
		logging.String("operation", "register"),
	)

	log := logging.L(ctx)

	req := &auth_v1.RegisterRequest{
		Username: user.Name,
		Password: user.Password,
	}

	log.Debug(
		"register handler called",
	)

	u, err := s.client.Register(ctx, req)
	if err != nil {
		log.Error(
			"failed to register",
			logging.Err(err),
		)

		return nil, err
	}

	log.Debug(
		"register success",
	)

	out := &models.User{
		ID:   uuid.MustParse(u.Id),
		Name: u.Username,
	}

	return out, nil
}

func (s *Service) Authenticate(ctx context.Context, accessToken string) (*models.User, error) {

	logging.WithAttrs(ctx,
		logging.String("operation", "authenticate"),
	)

	req := &auth_v1.AuthenticateRequest{
		AccessToken: accessToken,
	}

	log := logging.L(ctx)

	log.Debug(
		"authenticate handler called",
	)

	u, err := s.client.Authenticate(ctx, req)
	if err != nil {
		log.Error(
			"failed to authenticate",
			logging.Err(err),
		)

		return nil, err
	}

	log.Debug(
		"authenticate success",
		logging.String("id", u.Id),
		logging.String("username", u.Username),
	)

	out := &models.User{
		ID:   uuid.MustParse(u.Id),
		Name: u.Username,
	}

	return out, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*models.Tokens, error) {

	logging.WithAttrs(ctx,
		logging.String("operation", "refresh"),
	)

	log := logging.L(ctx)

	req := &auth_v1.RefreshRequest{
		RefreshToken: refreshToken,
	}

	log.Debug(
		"refresh handler called",
	)

	tokens, err := s.client.Refresh(ctx, req)
	if err != nil {
		log.Error(
			"failed to refresh",
			logging.Err(err),
		)
		return nil, err
	}

	log.Debug(
		"refresh success",
		logging.String("access_token", tokens.AccessToken),
		logging.String("refresh_token", tokens.RefreshToken),
	)

	out := &models.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	return out, nil
}
