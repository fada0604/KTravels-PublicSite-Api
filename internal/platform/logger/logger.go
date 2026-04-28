package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	log zerolog.Logger
}

func New(service, level, format string) (*Logger, error) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	if service == "" {
		service = "ktravels-publicsite-api"
	}

	var output io.Writer = os.Stdout
	if format == "console" {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}

	log := zerolog.New(output).
		Level(lvl).
		With().
		Str("service", service).
		Timestamp().
		Caller().
		Logger()

	return &Logger{log: log}, nil
}

func (l *Logger) Debug() *zerolog.Event {
	return l.log.Debug()
}

func (l *Logger) Info() *zerolog.Event {
	return l.log.Info()
}

func (l *Logger) Warn() *zerolog.Event {
	return l.log.Warn()
}

func (l *Logger) Error() *zerolog.Event {
	return l.log.Error()
}

func (l *Logger) With() zerolog.Context {
	return l.log.With()
}
