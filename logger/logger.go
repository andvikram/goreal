package logger

import (
	"fmt"
	"io"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/andvikram/goreal/configuration"
)

// Fields ...
type Fields map[string]interface{}

// GoRealLog ...
type GoRealLog struct{}

// GoRealLogging ...
type GoRealLogging interface {
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
	WithFields(f Fields) *GoRealLog
}

const (
	dirName    = "logs"
	logFileExt = ".log"
)

var (
	logger       *log.Logger
	logrusFields log.Fields
	logDir       string
	logFilePath  string
	logFile      *os.File
	env          string
	err          error
	once         sync.Once
	pathSep      = string(os.PathSeparator)
)

// Start will initialize a new service level logger
func Start(serviceEnv string) {
	env = serviceEnv
	once.Do(func() {
		new()
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
func (gorealLog *GoRealLog) Debug(message string) {
	logger.WithFields(logrusFields).Debug(message)
}

// Info ...
func (gorealLog *GoRealLog) Info(message string) {
	logger.WithFields(logrusFields).Info(message)
}

// Warn ...
func (gorealLog *GoRealLog) Warn(message string) {
	logger.WithFields(logrusFields).Warn(message)
}

// Error ...
func (gorealLog *GoRealLog) Error(message string) {
	logger.WithFields(logrusFields).Error(message)
}

// Fatal ...
func (gorealLog *GoRealLog) Fatal(message string) {
	logger.WithFields(logrusFields).Fatal(message)
}

// Panic ...
func (gorealLog *GoRealLog) Panic(message string) {
	logger.WithFields(logrusFields).Panic(message)
}

// WithFields ...
func (gorealLog *GoRealLog) WithFields(f Fields) *GoRealLog {
	logrusFields = make(log.Fields, len(f))
	for k, v := range f {
		logrusFields[k] = v
	}
	return gorealLog
}

func new() {
	config := configuration.Config
	logFile, err = os.OpenFile(config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Cannot open file, error:", err)
	}

	logger = log.New()

	if env == configuration.DevEnv {
		logger.Out = io.MultiWriter(os.Stdout, logFile)
	} else {
		logger.Out = logFile
	}

	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		level = log.InfoLevel
	}
	logger.Level = level

	timestampFormat := "2006-01-02 15:04:05.999Z07:00"
	logger.Formatter = &log.TextFormatter{
		TimestampFormat: timestampFormat,
	}
	if env == configuration.ProdEnv {
		logger.Formatter = &log.JSONFormatter{
			TimestampFormat: timestampFormat,
		}
	}
}
