package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

const (
	logDir  = "logs"
	infoLog = "info.log"
	errLog  = "error.log"
)

func Init() error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	infoWriter, err := os.OpenFile(
		filepath.Join(logDir, infoLog),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}

	errorWriter, err := os.OpenFile(
		filepath.Join(logDir, errLog),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), zapcore.InfoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), zapcore.ErrorLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Log = logger.Sugar()

	return nil
}
