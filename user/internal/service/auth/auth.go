package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/internal/models"
	"github.com/meraiku/micro/user/pkg/tokens"
)

var (
	accessSecret  = os.Getenv("ACCESS_SECRET")
	refreshSecret = os.Getenv("REFRESH_SECRET")

	accessTTL  = 24 * time.Hour
	refreshTTL = 7 * 24 * time.Hour
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInvalidTokens     = errors.New("invalid tokens")
	ErrNoTokens          = errors.New("no tokens")
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
}

type TokenRepository interface {
	StashTokens(ctx context.Context, userID string, tokens *models.Tokens) error
	GetTokens(ctx context.Context, userID string) (*models.Tokens, error)
}

type Service struct {
	userRepo  UserRepository
	tokenRepo TokenRepository

	accessTTL  time.Duration
	refreshTTL time.Duration

	accessSecret  []byte
	refreshSecret []byte
}

func New(
	userRepo UserRepository,
	tokenRepo TokenRepository,
) *Service {
	return &Service{
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (s *Service) Login(ctx context.Context, user *models.User) (*models.Tokens, error) {

	u, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := u.ValidatePassword(user.Password); err != nil {
		return nil, ErrIncorrectPassword
	}

	tokens, err := tokens.GeneratePair(
		u.ID.String(),
		s.accessTTL,
		s.refreshTTL,
		s.accessSecret,
		s.refreshSecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	if err := s.tokenRepo.StashTokens(ctx, user.ID.String(), tokens); err != nil {
		return nil, fmt.Errorf("failed to stash tokens: %w", err)
	}

	return tokens, nil
}

func (s *Service) Register(ctx context.Context, user *models.User) (*models.User, error) {

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	u, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return u, nil
}

func (s *Service) GetTokens(ctx context.Context, user *models.User) (*models.Tokens, error) {
	return nil, nil
}

func (s *Service) Authenticate(ctx context.Context, t *models.Tokens) (*models.User, error) {
	log := logging.L(ctx)

	log.Debug(
		"parsing access token",
		logging.String("token", t.AccessToken),
	)

	claims, err := tokens.ParseJWT(t.AccessToken, s.accessSecret)
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

	if repoTokens.AccessToken != t.AccessToken || repoTokens.RefreshToken != t.RefreshToken {
		return nil, ErrInvalidTokens
	}

	log.Debug(
		"get user from repo",
		logging.String("user_id", userID),
	)

	user, err := s.userRepo.GetByID(ctx, uuid.MustParse(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

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

	return tokens, nil
}
