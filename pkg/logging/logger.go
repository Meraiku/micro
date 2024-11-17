package logging

import (
	"context"
	"log/slog"
	"os"
)

const (
	defaultLevel      = LevelInfo
	defaultIsJSON     = true
	defaultAddSource  = true
	defaultSetDefault = true
	defaultLogstash   = false
)

func NewLogger(opts ...LoggerOption) *Logger {
	cfg := LoggerOptions{
		Level:      defaultLevel,
		IsJSON:     defaultIsJSON,
		AddSource:  defaultAddSource,
		SetDefault: defaultSetDefault,
		Logstash: Logstash{
			Enable: defaultLogstash,
		},
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	ho := &HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	var h Handler

	switch cfg.IsJSON {
	case true:
		h = NewJSONHandler(os.Stdout, ho)
	case false:
		h = NewTextHandler(os.Stdout, ho)
	}

	l := New(h)

	if cfg.SetDefault {
		SetDefault(l)
	}

	return l
}

type LoggerOptions struct {
	Level      Level
	IsJSON     bool
	AddSource  bool
	SetDefault bool
	Logstash   Logstash
}

type Logstash struct {
	Enable bool
	Addr   string
}

type LoggerOption func(*LoggerOptions)

func WithLevel(level Level) LoggerOption {
	return func(o *LoggerOptions) {
		var l Level

		if err := l.UnmarshalText([]byte(level.String())); err != nil {
			l = LevelInfo
		}

		o.Level = l
	}
}

func WithJSON(isJSON bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.IsJSON = isJSON
	}
}

func WithSource(addSource bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.AddSource = addSource
	}
}

func WithSetDefault(setDefault bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.SetDefault = setDefault
	}
}

func WithLogstash(logstashAddress string) LoggerOption {
	return func(o *LoggerOptions) {
		logstash := Logstash{
			Enable: true,
			Addr:   logstashAddress,
		}
		o.Logstash = logstash
	}
}

func WithAttrs(ctx context.Context, attrs ...Attr) *Logger {
	logger := L(ctx)

	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

func WithDefaultAttrs(logger *Logger, attrs ...Attr) *Logger {

	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

func L(ctx context.Context) *Logger {
	return loggerFromContext(ctx)
}

func Default() *Logger {
	return slog.Default()
}
