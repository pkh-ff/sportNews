package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
)

// InitLogger
// 初始化 Logger，可根據 `isDev` 參數決定是開發模式還是生產模式
func InitLogger(isDev bool) {
	var cfg zap.Config

	if isDev {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 彩色輸出
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 設定時間格式
		// cfg.OutputPaths 不需要設定，因為我們自己組 MultiWriteSyncer
	}

	// 檔案輪轉
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    50,
		MaxBackups: 5,
		MaxAge:     7,
		Compress:   true,
	})

	// log 等級
	level := zap.InfoLevel
	if isDev {
		level = zap.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), fileWriter),
		level,
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugaredLogger = logger.Sugar()
}

// CloseLogger
// 在程式結束時呼叫，確保日誌同步寫入
func CloseLogger() {
	_ = logger.Sync()
}

// Info uses fmt.Sprint to construct and log a message.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	sugaredLogger.Infof(template, args...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Infow(msg, keysAndValues...)
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	sugaredLogger.Debugf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//
//	s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Debugw(msg, keysAndValues...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	sugaredLogger.Warnf(template, args...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Warnw(msg, keysAndValues...)
}

// Error One uses fmt.Sprint to construct and log a message.
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	sugaredLogger.Errorf(template, args...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Errorw(msg, keysAndValues...)
}
