package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(lvl zapcore.Level, format string) (*zap.Logger, error) {
	cfg := defaultZapConfig()
	cfg.Level.SetLevel(lvl)
	cfg.Encoding = format

	return cfg.Build()
}

func defaultZapConfig() zap.Config {
	cfg := zap.NewProductionConfig()

	cfg.Encoding = "console"

	cfg.EncoderConfig.NameKey = "logger"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.CallerKey = "caller"
	// cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	cfg.OutputPaths = []string{"stderr"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	return cfg
}
