package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/meraiku/micro/user/internal/models"
)

var (
	ErrNoSecret = errors.New("empty jwt secret")
)

type Claims struct {
	ID       string `json:"id"`
	UID      string `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(
	id string,
	username string,
	ttl time.Duration,
	secret []byte,
) (string, error) {

	if len(secret) == 0 {
		return "", ErrNoSecret
	}

	c := &Claims{
		ID:       id,
		UID:      uuid.NewString(),
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl).UTC()),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	token, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GeneratePair(
	id string,
	username string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	accessSecret []byte,
	refreshSecret []byte,
) (*models.Tokens, error) {

	access, err := GenerateJWT(id, username, accessTTL, accessSecret)
	if err != nil {
		return nil, err
	}

	refresh, err := GenerateJWT(id, username, refreshTTL, refreshSecret)
	if err != nil {
		return nil, err
	}

	return models.NewTokens(access, refresh)
}

func ParseJWT(tokenStr string, secret []byte) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if token != nil {
		if token.Valid {
			return claims, nil
		}
	}

	return nil, err
}
