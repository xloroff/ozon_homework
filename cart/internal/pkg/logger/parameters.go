package logger

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"go.uber.org/zap/zapcore"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/config"
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

// LogWriter заводим кастомный записывальщик логов.
func LogWriter(folder string) (zapcore.WriteSyncer, error) {
	logFile, err := os.OpenFile(folder+config.LogfileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		// Если любая ошибка кроме того, что файлик просто не создан - прекращаем попытки.
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("Невозможно получить доступ к файлу логирования: %s - %w.", folder+config.LogfileName, err)
		}
		// Создаем если файла просто нет.
		logFile, err = os.Create(folder + config.LogfileName)
		if err != nil {
			return nil, fmt.Errorf("Невозможно создать файл логирования: %s - %w.", folder+config.LogfileName, err)
		}
	}

	// TODO добавить ротацию файлов логов через github.com/natefinch/lumberjack - лить все в один файлик не очень.

	return zapcore.AddSync(logFile), nil
}