package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

func New() error {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	infoFile, err := os.OpenFile("logs/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	errorFile, err := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l < zapcore.ErrorLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.ErrorLevel
	})

	infoCore := zapcore.NewCore(encoder, zapcore.AddSync(infoFile), infoLevel)
	errorCore := zapcore.NewCore(encoder, zapcore.AddSync(errorFile), errorLevel)

	core := zapcore.NewTee(infoCore, errorCore)

	Log = zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	Sugar = Log.Sugar()
	return nil
}
