package logging

import "os"

const (
	defaultLevel     = LevelInfo
	defaultIsJSON    = true
	defaultAddSource = true
	defaultLogstash  = false
)

func NewLogger(opts ...LoggerOption) *Logger {
	cfg := LoggerOptions{
		Level:     defaultLevel,
		IsJSON:    defaultIsJSON,
		AddSource: defaultAddSource,
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

	return l
}

type LoggerOptions struct {
	Level     Level
	IsJSON    bool
	AddSource bool
	Logstash  Logstash
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

func WithLogstash(logstashAddress string) LoggerOption {
	return func(o *LoggerOptions) {
		logstash := Logstash{
			Enable: true,
			Addr:   logstashAddress,
		}
		o.Logstash = logstash
	}
}
