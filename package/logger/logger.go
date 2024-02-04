package logger

import (
	"go.uber.org/zap"
)

func Channel(channel string) *zap.Logger {
	return loggers[channel]
}

func Debug(channel string, message string, fields ...zap.Field) {
	loggers[channel].Debug(message, fields...)
}

func Info(channel string, message string, fields ...zap.Field) {
	loggers[channel].Info(message, fields...)
}

func Warn(channel string, message string, fields ...zap.Field) {
	loggers[channel].Warn(message, fields...)
}

func Error(channel string, message string, fields ...zap.Field) {
	loggers[channel].Error(message, fields...)
}

func Fatal(channel string, message string, fields ...zap.Field) {
	loggers[channel].Fatal(message, fields...)
}

func Panic(channel string, message string, fields ...zap.Field) {
	loggers[channel].Panic(message, fields...)
}
