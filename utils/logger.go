package utils

import (
	"fmt"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	enabled Level
}

func New(level Level) *Logger {
	return &Logger{enabled: level}
}

func (l *Logger) log(level Level, color, tag string, msg string, args ...any) {
	if level < l.enabled {
		return
	}

	prefix := fmt.Sprintf("%s[Chip8]: %s [%s]",
		color, Reset, tag,
	)

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	fmt.Println(prefix, msg)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(DEBUG, Gray, "DEBUG", msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(INFO, Green, "INFO", msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(WARN, Yellow, "WARN", msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(ERROR, Red, "ERROR", msg, args...)
}

how to use it ? on any .go