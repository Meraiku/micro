package models

import "errors"

var (
	ErrEmptyAccessToken  = errors.New("empty access token")
	ErrEmptyRefreshToken = errors.New("empty refresh token")
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewTokens(accessToken, refreshToken string) (*Tokens, error) {
	if accessToken == "" {
		return nil, ErrEmptyAccessToken
	}

	if refreshToken == "" {
		return nil, ErrEmptyRefreshToken
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
