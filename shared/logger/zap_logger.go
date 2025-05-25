package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

var (
	instance *ZapLogger
	once     sync.Once
)

func InitLogger() (Logger, error) {
	cfg := zap.NewProductionConfig()

	cfg.Encoding = "console"

	// Cấu hình encoder cho console, bắt buộc set các encoder để tránh lỗi runtime
	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Cần có màu cho level
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 thời gian dễ đọc
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapLog, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.DPanicLevel))
	if err != nil {
		return nil, err
	}
	zapLog.Info("Zap logger initialized")
	instance = &ZapLogger{logger: zapLog}
	return instance, nil
}

func GetLogger() Logger {
	once.Do(func() {
		InitLogger()
	})
	return instance
}

func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, fields...)
}

func (z *ZapLogger) Error(msg string, fields ...Field) {
	z.logger.Error(msg, fields...)
}

func (z *ZapLogger) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, fields...)
}

func (z *ZapLogger) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, fields...)
}

func (z *ZapLogger) With(fields ...Field) Logger {
	return &ZapLogger{logger: z.logger.With(fields...)}
}
