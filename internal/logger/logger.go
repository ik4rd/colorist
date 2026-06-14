package logger

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int

const (
	LevelInfo Level = iota
	LevelWarn
	LevelError
	LevelPanic
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelPanic:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}

const (
	timeFormat    = "2006/01/02 15:04:05"
	fatalExitCode = 1
)

type Logger struct {
	mu       sync.Mutex
	out      io.Writer
	minLevel Level
}

func New(out io.Writer) *Logger {
	return &Logger{out: out, minLevel: LevelInfo}
}

func (l *Logger) SetMinLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.minLevel = level
}

func (l *Logger) write(level Level, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.minLevel {
		return
	}

	fmt.Fprintf(l.out, "%s %-5s %s\n", time.Now().Format(timeFormat), level, msg)
}

func (l *Logger) Infof(format string, args ...any) {
	l.write(LevelInfo, fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.write(LevelWarn, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(err error) {
	l.write(LevelError, err.Error())
}

func (l *Logger) Fatal(err error) {
	l.write(LevelError, err.Error())
	os.Exit(fatalExitCode)
}

func (l *Logger) Recover() {
	if r := recover(); r != nil {
		l.write(LevelPanic, fmt.Sprintf("%v\n%s", r, debug.Stack()))
		os.Exit(fatalExitCode)
	}
}
