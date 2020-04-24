package log

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type PgxAdapter struct {
	l Logger
}

func NewPgxAdapter(l Logger) *PgxAdapter {
	return &PgxAdapter{l: l}
}

func (l *PgxAdapter) Log(_ context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	var logger Logger
	if data != nil {
		logger = l.l.WithFields(data)
	} else {
		logger = l.l
	}

	switch level {
	case pgx.LogLevelTrace:
		logger.WithFields(map[string]interface{}{"PGX_LOG_LEVEL": level}).Debug(msg)
	case pgx.LogLevelDebug:
		logger.Debug(msg)
	case pgx.LogLevelInfo:
		logger.Info(msg)
	case pgx.LogLevelWarn:
		logger.Warn(msg)
	case pgx.LogLevelError:
		logger.Error(msg)
	default:
		logger.WithFields(map[string]interface{}{"PGX_LOG_LEVEL": level}).Error(msg)
	}
}
