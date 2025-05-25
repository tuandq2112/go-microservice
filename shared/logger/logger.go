package logger

import "go.uber.org/zap"

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	With(fields ...Field) Logger
}


type Field = zap.Field // Re-export để tránh import zap trong domain/usecase
