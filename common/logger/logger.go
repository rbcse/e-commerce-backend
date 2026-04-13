package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level string
}

var (
	once sync.Once
	lgr  *zap.Logger
)

func InitLogger(cfg Config) {
	once.Do(func() {
		lgr = newConsoleLogger(parseLevel(cfg.Level))
	})
}

func GetLogger() *zap.Logger {
	once.Do(func() {
		lgr = newConsoleLogger(zapcore.InfoLevel)
	})
	return lgr
}

func Debug(msg string, args ...interface{}) {
	GetLogger().Debug(fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...interface{}) {
	GetLogger().Info(fmt.Sprintf(msg, args...))
}

func Warn(msg string, args ...interface{}) {
	GetLogger().Warn(fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...interface{}) {
	GetLogger().Error(fmt.Sprintf(msg, args...))
}

func Fatal(msg string, args ...interface{}) {
	GetLogger().Fatal(fmt.Sprintf(msg, args...))
}

func newConsoleLogger(level zapcore.Level) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, 
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg), 
		zapcore.AddSync(coloredStdout()),      
		zap.NewAtomicLevelAt(level),
	)

	return zap.New(core,
		zap.AddCaller(),      
		zap.AddCallerSkip(1), 
	)
}

func coloredStdout() zapcore.WriteSyncer {
	return zapcore.Lock(zapcore.AddSync(os.Stdout))
}

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

func parseLevel(level string) zapcore.Level {
	if l, ok := levelMap[level]; ok {
		return l
	}
	return zapcore.InfoLevel 
}