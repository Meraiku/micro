package containers

import (
	"fmt"

	"github.com/meraiku/micro/user/internal/config"
	tokenRepo "github.com/meraiku/micro/user/internal/domain/token/memory"
	userRepo "github.com/meraiku/micro/user/internal/domain/user/memory"
	"github.com/meraiku/micro/user/internal/service/auth"
	"github.com/meraiku/micro/user/internal/service/user"
)

type AuthServiceRepos struct {
	user  auth.UserRepository
	token auth.TokenRepository
}

type UserServiceRepos struct {
	user user.Repository
}

func NewAuthServiceRepos(repos map[config.Repo]config.RepoType) (*AuthServiceRepos, error) {
	repo := &AuthServiceRepos{}

	user, ok := repos[config.Users]
	if !ok {
		return nil, fmt.Errorf("%w: %s", config.ErrRepoNotFound, "users")
	}

	token, ok := repos[config.Tokens]
	if !ok {
		return nil, fmt.Errorf("%w: %s", config.ErrRepoNotFound, "tokens")
	}

	switch user {
	case config.Memory:
		memoryUserRepo := userRepo.New()
		repo.user = memoryUserRepo
	}

	switch token {
	case config.Memory:
		memoryTokenRepo := tokenRepo.New()
		repo.token = memoryTokenRepo
	}

	return repo, nil

}

func NewUserServiceRepos(repos map[config.Repo]config.RepoType) (*UserServiceRepos, error) {

	repo := &UserServiceRepos{}

	user, ok := repos[config.Users]
	if !ok {
		return nil, fmt.Errorf("%w: %s", config.ErrRepoNotFound, "users")
	}

	switch user {
	case config.Memory:
		memoryUserRepo := userRepo.New()
		repo.user = memoryUserRepo
	}

	return repo, nil
}
