package logging

const (
	defaultLevel = LevelInfo
)

func New() *Logger {
	return &Logger{}
}

type LoggerOptions struct {
	Level Level
}
