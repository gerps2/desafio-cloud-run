package logger

import (
	"log"
	"os"
)

//go:generate mockery --name=Logger
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

type logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
}

func New() Logger {
	return &logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *logger) Info(msg string, args ...interface{}) {
	l.infoLogger.Printf(msg, args...)
}

func (l *logger) Error(msg string, args ...interface{}) {
	l.errorLogger.Printf(msg, args...)
}

func (l *logger) Debug(msg string, args ...interface{}) {
	l.debugLogger.Printf(msg, args...)
}

func (l *logger) Warn(msg string, args ...interface{}) {
	l.warnLogger.Printf(msg, args...)
}
