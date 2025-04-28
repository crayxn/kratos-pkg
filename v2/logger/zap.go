package logger

import (
	kz "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

const (
	STDOUT = 0
	FILE   = 1
	ALL    = 2
)

type ZapConfig struct {
	Level      string `json:"level;omitempty"`
	Writer     int32  `json:"writer;omitempty"`
	FileName   string `json:"file_name;omitempty"`
	MaxSize    int32  `json:"max_size;omitempty"`
	MaxBackups int32  `json:"max_backups;omitempty"`
	MaxAge     int32  `json:"max_age;omitempty"`
	Compress   bool   `json:"compress;omitempty"`
}

func NewZapLogger(conf *ZapConfig) log.Logger {
	// choose logger
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "module",
		TimeKey:        "time",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	var writer io.Writer
	if conf.Writer == ALL {
		writer = io.MultiWriter(getZapFileWriter(conf), os.Stdout)
	} else if conf.Writer == FILE {
		writer = getZapFileWriter(conf)
	} else {
		writer = os.Stdout
	}
	return kz.NewLogger(
		zap.New(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				zapcore.AddSync(writer),
				getZapLevel(conf.Level),
			),
			zap.WithCaller(false),
		),
	)
}

func getZapFileWriter(config *ZapConfig) io.Writer {
	// default
	if config.FileName == "" {
		config.FileName = "./runtime/log/app.log"
	}
	if config.MaxSize == 0 {
		config.MaxSize = 10
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 3
	}
	if config.MaxAge == 0 {
		config.MaxAge = 7
	}
	// lumberjack
	return &lumberjack.Logger{
		Filename:   config.FileName,
		MaxSize:    int(config.MaxSize),
		MaxBackups: int(config.MaxBackups),
		MaxAge:     int(config.MaxAge),
		Compress:   config.Compress,
	}
}

func getZapLevel(level string) zapcore.Level {
	switch log.ParseLevel(level) {
	case log.LevelDebug:
		return zapcore.DebugLevel
	case log.LevelInfo:
		return zapcore.InfoLevel
	case log.LevelWarn:
		return zapcore.WarnLevel
	case log.LevelError:
		return zapcore.ErrorLevel
	case log.LevelFatal:
		return zapcore.FatalLevel
	}
	return zapcore.InfoLevel
}
