package logger

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/andvikram/goreal/configuration"
)

// Fields type is used to pass to `WithFields`
type Fields map[string]interface{}

// GoRealLog ...
type GoRealLog interface {
	// Debug ...
	Debug(message string)
	// Debug ...
	Info(message string)
	// Debug ...
	Warn(message string)
	// Debug ...
	Error(message string)
	// Debug ...
	Fatal(message string)
	// Debug ...
	Panic(message string)
	// WithFields allows to add fields to logging
	WithFields(Fields) GoRealLog
}

type grLog struct {
	logger *logrus.Logger
	fields logrus.Fields
}

const (
	dirName    = "logs"
	logFileExt = ".log"
)

var (
	// Log is instance for GoRealLog interface
	Log         GoRealLog
	logDir      string
	logFilePath string
	logFile     *os.File
	env         string
	err         error
	once        sync.Once
	pathSep     = string(os.PathSeparator)
)

// Start will initialize a new service level logger
func Start(serviceEnv string) {
	env = serviceEnv
	once.Do(func() {
		Log = newLog()
	})
}

// Stop the logger
func Stop() error {
	err := logFile.Close()
	if err != nil {
		fmt.Println("Error closing log file, error:", err)
	}
	return err
}

// Debug ...
func (l *grLog) Debug(message string) {
	l.logger.WithFields(l.fields).Debug(message)
}

// Info ...
func (l *grLog) Info(message string) {
	l.logger.WithFields(l.fields).Info(message)
}

// Warn ...
func (l *grLog) Warn(message string) {
	l.logger.WithFields(l.fields).Warn(message)
}

// Error ...
func (l *grLog) Error(message string) {
	l.logger.WithFields(l.fields).Error(message)
}

// Fatal ...
func (l *grLog) Fatal(message string) {
	l.logger.WithFields(l.fields).Fatal(message)
}

// Panic ...
func (l *grLog) Panic(message string) {
	l.logger.WithFields(l.fields).Panic(message)
}

// WithFields ...
func (l *grLog) WithFields(f Fields) GoRealLog {
	l.fields = make(logrus.Fields, len(f))
	for k, v := range f {
		l.fields[k] = v
	}
	return l
}

func newLog() GoRealLog {
	config := configuration.Config
	logFile, err = os.OpenFile(config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		logrus.Fatal("Cannot open file, error:", err)
	}

	l := new(grLog)

	l.logger = logrus.New()

	if env == configuration.DevEnv {
		l.logger.Out = io.MultiWriter(os.Stdout, logFile)
	} else {
		l.logger.Out = logFile
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	l.logger.Level = level

	timestampFormat := "2006-01-02 15:04:05.999Z07:00"
	l.logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: timestampFormat,
	}
	if env == configuration.ProdEnv {
		l.logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: timestampFormat,
		}
	}

	return l
}
