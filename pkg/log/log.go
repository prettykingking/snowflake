package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/prettykingking/snowflake/pkg/config"
)

var zl *zap.Logger

func NewLogger(cf *config.Logging) (*zap.Logger, error) {
	logger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(newLevel(cf.Level)),
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()

	if err != nil {
		return nil, err
	}

	zl = logger

	return logger, err
}

func GetLogger() *zap.Logger {
	return zl
}

// newLevel returns logging level based on string either in UPPER case or lower case
func newLevel(text string) zapcore.Level {
	var l zapcore.Level
	_ = l.Set(text) // defaults to INFO level
	return l
}
