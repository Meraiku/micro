package config

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/meraiku/micro/pkg/logging"
)

var (
	ErrNoSecret = errors.New("empty jwt secret")
)

type cfgAPI interface {
	Address() string
}

type API byte

const (
	GRPC API = iota + 1
	REST
)

type Config struct {
	apiType API
	Api     cfgAPI
	Secrets JWTSecrets
	TTL     TokenTTL
	Repos   map[Repo]RepoType
	Brokers []string
	Topic   string
}

func NewConfig(ctx context.Context, api API) (*Config, error) {

	c := &Config{
		apiType: api,
	}

	if err := c.initConfig(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Address() string {
	return c.Api.Address()
}

func (c *Config) initConfig(ctx context.Context) error {

	deps := []func(ctx context.Context) error{
		c.setupAPI,
		c.setupJWT,
		c.setupRepos,
		c.setupBrokers,
	}

	for _, dep := range deps {
		if err := dep(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) setupAPI(ctx context.Context) error {

	switch c.apiType {
	case REST:
		a, err := NewREST()
		if err != nil {
			return err
		}
		c.Api = a
	case GRPC:
		a, err := NewGRPC()
		if err != nil {
			return err
		}
		c.Api = a
	}

	return nil
}

func (c *Config) setupJWT(ctx context.Context) error {

	tokenTTL, err := NewTTL()
	if err != nil {
		return err
	}

	c.TTL = tokenTTL

	jwtSecrets, err := NewSecrets()
	if err != nil {
		return err
	}

	c.Secrets = jwtSecrets

	logging.L(ctx).Info("jwt initialized")

	return nil
}

func (c *Config) setupRepos(ctx context.Context) error {
	userType := os.Getenv("USER_REPO")
	if userType == "" {
		userType = "memory"
	}

	tokensType := os.Getenv("TOKENS_REPO")
	if tokensType == "" {
		tokensType = "memory"
	}

	c.Repos = map[Repo]RepoType{
		Users:  ConvertToType(userType),
		Tokens: ConvertToType(tokensType),
	}

	return nil
}

func (c *Config) setupBrokers(ctx context.Context) error {
	brokers := os.Getenv("KAFKA_BROKERS")
	c.Brokers = strings.Split(brokers, ",")

	c.Topic = os.Getenv("KAFKA_TOPIC")
	if c.Topic == "" {
		c.Topic = "user"
	}

	logging.L(ctx).Info("brokers initialized", logging.Any("brokers", c.Brokers), logging.String("topic", c.Topic))

	return nil
}

type TokenTTL struct {
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewTTL() (TokenTTL, error) {
	accessTTL := os.Getenv("ACCESS_TTL")
	refreshTTL := os.Getenv("REFRESH_TTL")

	if accessTTL == "" {
		accessTTL = "24h"
	}

	if refreshTTL == "" {
		// 7 days
		refreshTTL = "168h"
	}

	attl, err := time.ParseDuration(accessTTL)
	if err != nil {
		return TokenTTL{}, err
	}

	rttl, err := time.ParseDuration(refreshTTL)
	if err != nil {
		return TokenTTL{}, err
	}

	return TokenTTL{
		AccessTTL:  attl,
		RefreshTTL: rttl,
	}, nil
}

type JWTSecrets struct {
	AccessSecret  string
	RefreshSecret string
}

func NewSecrets() (JWTSecrets, error) {
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")

	if accessSecret == "" || refreshSecret == "" {
		return JWTSecrets{}, ErrNoSecret
	}
	return JWTSecrets{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
	}, nil
}

func Load() {

	godotenv.Load(".env")

}
