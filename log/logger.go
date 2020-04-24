package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})

	Println(v ...interface{})
	Printf(format string, v ...interface{})

	WithFields(fields map[string]interface{}) Logger
}

type LoggerLogrus struct {
	logrus.FieldLogger
	file *os.File
}

func New(logger *logrus.Logger, file *os.File) *LoggerLogrus {
	return &LoggerLogrus{
		FieldLogger: logger,
		file:        file,
	}
}

func (l *LoggerLogrus) WithFields(fields map[string]interface{}) Logger {
	return &LoggerLogrus{
		FieldLogger: l.FieldLogger.WithFields(fields),
	}
}

func (l *LoggerLogrus) Close() {
	if l.file == nil {
		return
	}
	err := l.file.Close()
	if err != nil {
		l.Errorf("Error closing log file: %s\n", err)
	}
}
