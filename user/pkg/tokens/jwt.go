package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrNoSecret = errors.New("empty jwt secret")
)

type Claims struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

func GenerateJWT(
	id string,
	ttl time.Duration,
	secret []byte,
) (string, error) {
	if secret == "" {
		return "", ErrNoSecret
	}

	c := &Claims{
		ID:  id,
		UID: uuid.NewString(),
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
