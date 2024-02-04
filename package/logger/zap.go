package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Initialize(config Config) error {
	loggers = map[string]*zap.Logger{}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	for _, channel := range config.Channels {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   channel.Filename,
			MaxSize:    channel.MaxSize,    // megabytes
			MaxBackups: channel.MaxBackups, // number of backups
			MaxAge:     channel.MaxAge,     // days
			Compress:   channel.Compress,
		}

		writeSyncer := zapcore.AddSync(lumberJackLogger)

		var level zapcore.Level

		err := level.UnmarshalText([]byte(channel.Level))

		if err != nil {
			return err
		}

		core := zapcore.NewCore(jsonEncoder, writeSyncer, level)

		loggers[channel.Name] = zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	}

	return nil
}
