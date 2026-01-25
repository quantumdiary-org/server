package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	logger zerolog.Logger
}


func New(logLevel string, logFile string) (*Logger, error) {
	var writer io.Writer

	
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	if logFile != "" {
		
		rl, err := rotatelogs.New(
			logFile+".%Y%m%d",
			rotatelogs.WithLinkName(logFile),
			rotatelogs.WithMaxAge(7*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create rotating log file: %w", err)
		}
		writer = rl
	} else {
		writer = os.Stdout
	}

	logger := zerolog.New(writer).Level(level).With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{logger: logger}, nil
}


func Init(logLevel string, logFile string) error {
	logger, err := New(logLevel, logFile)
	if err != nil {
		return err
	}

	log.Logger = logger.logger
	return nil
}


func (l *Logger) Info(message string, fields ...interface{}) {
	event := l.logger.Info()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				event.Str(key, fmt.Sprintf("%v", fields[i+1]))
			}
		}
	}
	event.Msg(message)
}


func (l *Logger) Error(message string, fields ...interface{}) {
	event := l.logger.Error()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				event.Str(key, fmt.Sprintf("%v", fields[i+1]))
			}
		}
	}
	event.Msg(message)
}


func (l *Logger) Debug(message string, fields ...interface{}) {
	event := l.logger.Debug()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				event.Str(key, fmt.Sprintf("%v", fields[i+1]))
			}
		}
	}
	event.Msg(message)
}


func (l *Logger) Warn(message string, fields ...interface{}) {
	event := l.logger.Warn()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				event.Str(key, fmt.Sprintf("%v", fields[i+1]))
			}
		}
	}
	event.Msg(message)
}


func (l *Logger) Fatal(message string, fields ...interface{}) {
	event := l.logger.Fatal()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				event.Str(key, fmt.Sprintf("%v", fields[i+1]))
			}
		}
	}
	event.Msg(message)
}


func (l *Logger) With(fields ...interface{}) *Logger {
	newLogger := l.logger.With()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				newLogger = newLogger.Str(key, fmt.Sprintf("%v", fields[i+1]))
			}
		}
	}
	return &Logger{logger: newLogger.Logger()}
}


func Info(message string, fields ...interface{}) {
	log.Info().Fields(getFieldMap(fields)).Msg(message)
}

func Error(message string, fields ...interface{}) {
	log.Error().Fields(getFieldMap(fields)).Msg(message)
}

func Debug(message string, fields ...interface{}) {
	log.Debug().Fields(getFieldMap(fields)).Msg(message)
}

func Warn(message string, fields ...interface{}) {
	log.Warn().Fields(getFieldMap(fields)).Msg(message)
}

func Fatal(message string, fields ...interface{}) {
	log.Fatal().Fields(getFieldMap(fields)).Msg(message)
}

func getFieldMap(fields []interface{}) map[string]interface{} {
	fieldMap := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				fieldMap[key] = fields[i+1]
			}
		}
	}
	return fieldMap
}