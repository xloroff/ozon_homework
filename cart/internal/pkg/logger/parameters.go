package logger

import (
	"go.uber.org/zap/zapcore"
)

// DebugLevel определяем уровень дебага по параметру из конфигурации.
func DebugLevel(logLevel string) zapcore.Level {
	var level zapcore.Level

	switch logLevel {
	case "Debug":
		level = zapcore.DebugLevel
	case "Info":
		level = zapcore.InfoLevel
	case "Warn":
		level = zapcore.WarnLevel
	case "Error":
		level = zapcore.ErrorLevel
	case "DPanic":
		level = zapcore.DPanicLevel
	case "Panic":
		level = zapcore.PanicLevel
	case "Fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.DebugLevel
	}

	return level
}
