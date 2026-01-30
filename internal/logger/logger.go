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

	infoCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(infoWriter),
		zapcore.InfoLevel,
	).With([]zapcore.Field{})
	infoCore = &filteringCore{
		core: infoCore,
		level: func(level zapcore.Level) bool {
			return level < zapcore.ErrorLevel
		},
	}

	errorCore := zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), zapcore.ErrorLevel)

	core := zapcore.NewTee(infoCore, errorCore)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Log = logger.Sugar()

	return nil
}

type filteringCore struct {
	core  zapcore.Core
	level func(zapcore.Level) bool
}

func (fc *filteringCore) Enabled(level zapcore.Level) bool {
	return fc.level(level) && fc.core.Enabled(level)
}

func (fc *filteringCore) With(fields []zapcore.Field) zapcore.Core {
	return &filteringCore{
		core:  fc.core.With(fields),
		level: fc.level,
	}
}

func (fc *filteringCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !fc.Enabled(entry.Level) {
		return checked
	}
	return fc.core.Check(entry, checked)
}

func (fc *filteringCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if !fc.level(entry.Level) {
		return nil
	}
	return fc.core.Write(entry, fields)
}

func (fc *filteringCore) Sync() error {
	return fc.core.Sync()
}
