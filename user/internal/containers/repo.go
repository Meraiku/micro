package containers

import (
	"os"

	"github.com/meraiku/micro/pkg/logging"
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

func NewAuthServiceRepos() (*AuthServiceRepos, error) {
	var (
		tokenRepository auth.TokenRepository
		userRepository  auth.UserRepository
	)

	authRepoEnv := os.Getenv("AUTH_REPO")
	if authRepoEnv == "" {
		authRepoEnv = "memory"
	}

	userRepoEnv := os.Getenv("USER_REPO")
	if userRepoEnv == "" {
		userRepoEnv = "memory"
	}

	switch authRepoEnv {
	case "memory":
		memoryTokenRepo := tokenRepo.New()
		tokenRepository = memoryTokenRepo
	}

	switch userRepoEnv {
	case "memory":
		memoryUserRepo := userRepo.New()
		userRepository = memoryUserRepo
	}

	logging.Default().Info(
		"auth service initialized",
		logging.String("auth_repo", authRepoEnv),
		logging.String("user_repo", userRepoEnv),
	)

	repos := &AuthServiceRepos{
		user:  userRepository,
		token: tokenRepository,
	}

	return repos, nil

}

func NewUserServiceRepos() (*UserServiceRepos, error) {
	var (
		userRepository user.Repository
	)

	userRepoEnv := os.Getenv("USER_REPO")
	if userRepoEnv == "" {
		userRepoEnv = "memory"
	}

	switch userRepoEnv {
	case "memory":
		memoryUserRepo := userRepo.New()
		userRepository = memoryUserRepo
	}

	logging.Default().Info(
		"user service initialized",
		logging.String("user_repo", userRepoEnv),
	)

	repos := &UserServiceRepos{
		user: userRepository,
	}

	return repos, nil
}
