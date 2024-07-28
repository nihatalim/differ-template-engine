package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const timeFormat string = "2006-01-02T15:04:05.999Z"

type logger struct {
	*zap.SugaredLogger
}

func New(level zapcore.Level) Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "log",
		CallerKey:      "sourceLocation",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(timeFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	l, _ := config.Build(zap.AddStacktrace(zap.FatalLevel))
	return &logger{l.Sugar()}
}

func (l *logger) With(args ...interface{}) Logger {
	if len(args) > 0 {
		return &logger{l.SugaredLogger.With(args...)}
	}
	return l
}

func (l *logger) GetDesugaredLogger() *zap.Logger {
	return l.SugaredLogger.Desugar()
}
