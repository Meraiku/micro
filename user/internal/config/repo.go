package config

import "errors"

var (
	ErrRepoNotFound = errors.New("repo not found")
)

type RepoType byte

const (
	Memory RepoType = iota + 1
)

type Repo byte

const (
	Users Repo = iota + 1
	Tokens
)
