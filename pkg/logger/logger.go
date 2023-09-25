package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go-api/pkg/config"
)

type Fields map[any]any

// Logger is the interface that defines the main logger functions.
type Logger interface {
	Debug(message string, fields ...Fields)
	Info(message string, fields ...Fields)
	Warn(message string, fields ...Fields)
	Error(message string, fields ...Fields)
}

// NewZapLogger creates a new instance of Logger with the provided options.
func NewZapLogger(cfg *config.Config, options ...zap.Option) (Logger, error) {
	var level zapcore.Level
	if err := level.Set(cfg.Logger.Level); err != nil {
		return nil, err
	}

	config := zap.Config{
		Encoding:          "json",
		Level:             zap.NewAtomicLevelAt(level),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		DisableStacktrace: cfg.Logger.DisableStacktrace,
		DisableCaller:     cfg.Logger.DisableCaller,
		Development:       cfg.Logger.Development,
	}

	options = append(options, zap.AddCallerSkip(1))
	logger, err := config.Build(options...)
	if err != nil {
		return nil, err
	}

	return &zapLogger{logger.Sugar()}, nil
}

type zapLogger struct {
	sugarLogger *zap.SugaredLogger
}

func (l *zapLogger) Debug(message string, fields ...Fields) {
	l.sugarLogger.Debugw(message, toSlice(fields)...)
}

func (l *zapLogger) Info(message string, fields ...Fields) {
	l.sugarLogger.Infow(message, toSlice(fields)...)
}

func (l *zapLogger) Warn(message string, fields ...Fields) {
	l.sugarLogger.Warnw(message, toSlice(fields)...)
}

func (l *zapLogger) Error(message string, fields ...Fields) {
	l.sugarLogger.Errorw(message, toSlice(fields)...)
}

func (l *zapLogger) Fatal(message string, fields ...Fields) {
	l.sugarLogger.Fatalw(message, toSlice(fields)...)
}

func toSlice(fieldsArr []Fields) []any {
	fields := []any{}
	for _, f := range fieldsArr {
		fields = append(fields, f.toSlice()...)
	}
	return fields
}

func (f *Fields) toSlice() []any {
	if f == nil {
		return []any{}
	}

	keysAndValues := make([]any, len(*f)*2)
	i := 0
	for k, v := range *f {
		keysAndValues[i] = k
		keysAndValues[i+1] = v
		i += 2
	}
	return keysAndValues
}
