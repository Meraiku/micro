package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/internal/config"
	"github.com/meraiku/micro/user/internal/models"
	"github.com/meraiku/micro/user/pkg/kafka/producer"
	"github.com/meraiku/micro/user/pkg/tokens"
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInvalidTokens     = errors.New("invalid tokens")
	ErrNoTokens          = errors.New("no tokens")
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByUsername(ctx context.Context, name string) (*models.User, error)
}

type TokenRepository interface {
	StashTokens(ctx context.Context, userID string, tokens *models.Tokens) error
	GetTokens(ctx context.Context, userID string) (*models.Tokens, error)
}

type Service struct {
	userRepo  UserRepository
	tokenRepo TokenRepository
	notify    *producer.Producer

	accessTTL  time.Duration
	refreshTTL time.Duration

	accessSecret  []byte
	refreshSecret []byte
}

func New(
	cfg *config.Config,
	userRepo UserRepository,
	tokenRepo TokenRepository,
) (*Service, error) {

	notifier, err := producer.New(cfg.Brokers, cfg.Topic)
	if err != nil {
		return nil, err
	}

	notifier.Run()

	return &Service{
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		notify:        notifier,
		accessTTL:     cfg.TTL.AccessTTL,
		refreshTTL:    cfg.TTL.RefreshTTL,
		accessSecret:  []byte(cfg.Secrets.AccessSecret),
		refreshSecret: []byte(cfg.Secrets.RefreshSecret),
	}, nil
}

func (s *Service) Login(ctx context.Context, user *models.User) (*models.Tokens, error) {

	logging.WithAttrs(
		ctx,
		logging.String("username", user.Name),
	)

	log := logging.L(ctx)

	log.Debug(
		"get user from repo",
	)

	u, err := s.userRepo.GetByUsername(ctx, user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	log.Debug(
		"user from repo",
		logging.String("user_id", u.ID.String()),
		logging.String("username", u.Name),
	)

	log.Debug(
		"validate password",
	)

	if err := u.ValidatePassword(user.Password); err != nil {
		return nil, ErrIncorrectPassword
	}

	log.Debug(
		"generate tokens",
	)

	tokens, err := tokens.GeneratePair(
		u.ID.String(),
		u.Name,
		s.accessTTL,
		s.refreshTTL,
		s.accessSecret,
		s.refreshSecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	log.Debug(
		"stash tokens",
	)

	if err := s.tokenRepo.StashTokens(ctx, u.ID.String(), tokens); err != nil {
		return nil, fmt.Errorf("failed to stash tokens: %w", err)
	}

	log.Debug(
		"send notification",
	)

	go s.notify.Send(u.ID.String(), fmt.Sprintf("%s logged in", u.Name))

	return tokens, nil
}

func (s *Service) Register(ctx context.Context, user *models.User) (*models.User, error) {

	logging.WithAttrs(
		ctx,
		logging.String("username", user.Name),
	)

	log := logging.L(ctx)

	log.Debug(
		"creating user hash password",
	)

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	log.Debug(
		"create user in repo",
	)

	u, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Debug(
		"send notification",
	)

	go s.notify.Send(u.ID.String(), fmt.Sprintf("%s registered", u.Name))

	return u, nil
}

func (s *Service) GetTokens(ctx context.Context, user *models.User) (*models.Tokens, error) {
	return nil, nil
}

func (s *Service) Authenticate(ctx context.Context, accessToken string) (*models.User, error) {
	log := logging.L(ctx)

	log.Debug(
		"parsing access token",
		logging.String("token", accessToken),
	)

	claims, err := tokens.ParseJWT(accessToken, s.accessSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	userID := claims.ID

	log.Debug(
		"get tokens from repo",
		logging.String("user_id", userID),
	)

	repoTokens, err := s.tokenRepo.GetTokens(ctx, userID)
	if err != nil {
		return nil, ErrNoTokens
	}

	log.Debug(
		"tokens from repo",
		logging.String("access_token", repoTokens.AccessToken),
		logging.String("refresh_token", repoTokens.RefreshToken),
	)

	log.Debug(
		"compare tokens",
	)

	if repoTokens.AccessToken != accessToken {
		return nil, ErrInvalidTokens
	}

	user := &models.User{
		ID:   uuid.MustParse(userID),
		Name: claims.Username,
	}

	go s.notify.Send(userID, fmt.Sprintf("%s authenticated", userID))

	return user, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*models.Tokens, error) {
	log := logging.L(ctx)

	log.Debug(
		"parsing refresh token",
		logging.String("token", refreshToken),
	)

	claims, err := tokens.ParseJWT(refreshToken, s.refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	userID := claims.ID

	log.Debug(
		"get tokens from repo",
		logging.String("user_id", userID),
	)

	repoTokens, err := s.tokenRepo.GetTokens(ctx, userID)
	if err != nil {
		return nil, ErrNoTokens
	}

	log.Debug(
		"tokens from repo",
		logging.String("access_token", repoTokens.AccessToken),
		logging.String("refresh_token", repoTokens.RefreshToken),
	)

	if repoTokens.RefreshToken != refreshToken {
		return nil, ErrInvalidTokens
	}

	log.Debug(
		"generate new tokens",
	)

	tokens, err := tokens.GeneratePair(
		userID,
		claims.Username,
		s.accessTTL,
		s.refreshTTL,
		s.accessSecret,
		s.refreshSecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	log.Debug(
		"stash new tokens",
	)

	if err := s.tokenRepo.StashTokens(ctx, userID, tokens); err != nil {
		return nil, fmt.Errorf("failed to stash tokens: %w", err)
	}

	go s.notify.Send(userID, fmt.Sprintf("%s refreshed tokens", userID))

	return tokens, nil
}
