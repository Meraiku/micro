package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/internal/models"
	"github.com/meraiku/micro/user/pkg/auth_v1"
)

type AuthService interface {
	Login(ctx context.Context, user *models.User) (*models.Tokens, error)
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Authenticate(ctx context.Context, accessToken string) (*models.User, error)
	Refresh(ctx context.Context, refreshToken string) (*models.Tokens, error)
}

type GRPCAuthService struct {
	auth_v1.UnimplementedAuthV1Server

	authService AuthService
}

func NewAuthService(authService AuthService) *GRPCAuthService {
	return &GRPCAuthService{
		authService: authService,
	}
}

func (s *GRPCAuthService) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.Tokens, error) {
	logging.WithAttrs(
		ctx,
		logging.String("request_id", uuid.New().String()),
		logging.String("operation", "login"),
	)

	log := logging.L(ctx)

	user, err := models.NewUser(req.Username, req.Password)
	if err != nil {
		log.Error(
			"failed to create user",
			logging.Err(err),
		)
		return nil, err
	}

	tokens, err := s.authService.Login(ctx, user)
	if err != nil {
		log.Error(
			"failed to login user",
			logging.Err(err),
		)
		return nil, err
	}

	return FromTokens(tokens), nil
}

func (s *GRPCAuthService) Register(ctx context.Context, req *auth_v1.RegisterRequest) (*auth_v1.RegisterResponse, error) {

	logging.WithAttrs(
		ctx,
		logging.String("request_id", uuid.New().String()),
		logging.String("operation", "register"),
	)

	log := logging.L(ctx)

	log.Info(
		"register handler called",
		logging.String("username", req.Username),
	)

	user, err := models.NewUser(req.Username, req.Password)
	if err != nil {
		log.Error(
			"failed to create user",
			logging.Err(err),
		)
		return nil, err
	}

	log.Debug(
		"creating user",
		logging.String("username", user.Name),
	)

	user, err = s.authService.Register(ctx, user)
	if err != nil {
		log.Error(
			"failed to register user",
			logging.Err(err),
		)
		return nil, err
	}

	log.Debug(
		"user created",
		logging.String("id", user.ID.String()),
		logging.String("username", user.Name),
	)

	out := &auth_v1.RegisterResponse{
		Id:       user.ID.String(),
		Username: user.Name,
	}

	log.Info(
		"register handler done",
		logging.String("username", user.Name),
	)

	return out, nil
}

func (s *GRPCAuthService) Authenticate(ctx context.Context, req *auth_v1.AuthenticateRequest) (*auth_v1.User, error) {

	logging.WithAttrs(
		ctx,
		logging.String("request_id", uuid.New().String()),
		logging.String("operation", "authenticate"),
	)

	log := logging.L(ctx)

	accessToken := req.AccessToken

	user, err := s.authService.Authenticate(ctx, accessToken)
	if err != nil {
		log.Error(
			"failed to authenticate user",
			logging.Err(err),
		)
		return nil, err
	}

	out := &auth_v1.User{
		Id:       user.ID.String(),
		Username: user.Name,
	}

	return out, nil
}

func (s *GRPCAuthService) Refresh(ctx context.Context, req *auth_v1.RefreshRequest) (*auth_v1.Tokens, error) {

	logging.WithAttrs(
		ctx,
		logging.String("request_id", uuid.New().String()),
		logging.String("operation", "refresh"),
	)

	log := logging.L(ctx)

	refreshToken := req.RefreshToken

	tokens, err := s.authService.Refresh(ctx, refreshToken)
	if err != nil {
		log.Error(
			"failed to refresh tokens",
			logging.Err(err),
		)
		return nil, err
	}

	return FromTokens(tokens), nil
}
