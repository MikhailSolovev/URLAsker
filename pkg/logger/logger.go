package logger

import (
	"github.com/rs/zerolog"
	"os"
	"sync"
	"time"
)

// singleton logger

var Log Logger
var once sync.Once

type Logger struct {
	zerolog.Logger
}

func (l *Logger) LogDebug(msg string) {
	l.Debug().Msg(msg)
}

func (l *Logger) LogInfo(msg string) {
	l.Info().Msg(msg)
}

func (l *Logger) LogWarn(msg string) {
	l.Warn().Msg(msg)
}

func (l *Logger) LogError(msg string) {
	l.Error().Msg(msg)
}

func (l *Logger) LogFatal(msg string) {
	l.Fatal().Msg(msg)
}

func getLog(serviceName string, isPretty bool, logLevel string) zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Str("module", serviceName).Timestamp()

	if isPretty {
		logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Str("module", serviceName).Timestamp()
	}

	if logLevel != "" {
		level, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			// default log level is debug
			return logger.Logger()
		}
		zerolog.SetGlobalLevel(level)
	}

	return logger.Logger()
}

func New(serviceName string, isPretty bool, logLevel string) *Logger {
	once.Do(func() {
		Log = Logger{getLog(serviceName, isPretty, logLevel)}
	})

	return &Log
}
