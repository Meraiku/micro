package logging

import "log/slog"

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type (
	Logger         = slog.Logger
	Level          = slog.Level
	Handler        = slog.Handler
	HandlerOptions = slog.HandlerOptions
)

var (
	New            = slog.New
	NewTextHandler = slog.NewTextHandler
	NewJSONHandler = slog.NewJSONHandler
)
