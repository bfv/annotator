package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelNone  LogLevel = "none"
	LogLevelError LogLevel = "error"
	LogLevelInfo  LogLevel = "info"
	LogLevelDebug LogLevel = "debug"
	LogLevelTrace LogLevel = "trace"
)

// InitLogger initializes the logger based on configuration
func InitLogger(level LogLevel, logToConsole bool, logFile string) error {
	var writers []io.Writer

	// Determine log level
	var zeroLevel zerolog.Level
	switch level {
	case LogLevelNone:
		zeroLevel = zerolog.Disabled
	case LogLevelError:
		zeroLevel = zerolog.ErrorLevel
	case LogLevelInfo:
		zeroLevel = zerolog.InfoLevel
	case LogLevelDebug:
		zeroLevel = zerolog.DebugLevel
	case LogLevelTrace:
		zeroLevel = zerolog.TraceLevel
	default:
		zeroLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(zeroLevel)

	// Setup output writers
	if logToConsole {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		// Log to file
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		writers = append(writers, file)
	}

	if len(writers) == 0 {
		logger = zerolog.Nop()
	} else if len(writers) == 1 {
		logger = zerolog.New(writers[0]).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(io.MultiWriter(writers...)).With().Timestamp().Logger()
	}

	log.Logger = logger
	return nil
}
