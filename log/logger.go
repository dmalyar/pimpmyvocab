package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"reflect"
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

	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

type LoggerLogrus struct {
	logrus.FieldLogger
	out io.Closer
}

func New(logger *logrus.Logger, out io.Closer) *LoggerLogrus {
	return &LoggerLogrus{
		FieldLogger: logger,
		out:         out,
	}
}

func (l *LoggerLogrus) WithField(key string, value interface{}) Logger {
	return l.WithFields(map[string]interface{}{key: value})
}

func (l *LoggerLogrus) WithFields(fields map[string]interface{}) Logger {
	return &LoggerLogrus{
		FieldLogger: l.FieldLogger.WithFields(fields),
	}
}

func (l *LoggerLogrus) Close() {
	if l.out == nil || (reflect.ValueOf(l.out).Kind() == reflect.Ptr && reflect.ValueOf(l.out).IsNil()) {
		return
	}
	err := l.out.Close()
	if err != nil {
		l.Errorf("Error closing log file: %s", err)
	}
}
