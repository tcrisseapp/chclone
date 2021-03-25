package zaplogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New will initiate a new zap logger with the proper configuration
func New(logLevel string) (*zap.Logger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, err
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return zapConfig.Build()

}
