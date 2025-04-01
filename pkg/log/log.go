package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
)

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	sugaredLogger = logger.Sugar()
}

// InitLogger
// 初始化 Logger，可根據 `isDev` 參數決定是開發模式還是生產模式
func InitLogger(isDev bool) {
	var err error
	var cfg zap.Config

	if isDev {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 彩色輸出
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 設定時間格式
		cfg.OutputPaths = []string{"stdout", "logs/app.log"}      // 輸出到終端和檔案
	}

	logger, err = cfg.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log", // 日誌檔案路徑
		MaxSize:    50,             // 單檔最大 50MB
		MaxBackups: 5,              // 最多保留 5 個備份
		MaxAge:     7,              // 最多保留 7 天
		Compress:   true,           // 是否壓縮舊日誌
	})

	// 設定 log 等級
	level := zap.InfoLevel
	if isDev {
		level = zap.DebugLevel
	}

	// 建立 core，將 log 寫入終端與檔案
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),                           // JSON 格式
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), fileWriter), // 同時輸出到終端與檔案
		level, // 設定 log 等級
	)

	// 生成 logger
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugaredLogger = logger.Sugar()
}

// CloseLogger
// 在程式結束時呼叫，確保日誌同步寫入
func CloseLogger() {
	_ = logger.Sync()
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error One uses fmt.Sprint to construct and log a message.
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Panicw(msg, keysAndValues...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	sugaredLogger.DPanicf(template, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	sugaredLogger.Fatalf(template, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Fatalw(msg, keysAndValues...)
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

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	sugaredLogger.Errorf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	sugaredLogger.Panicf(template, args...)
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

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	sugaredLogger.Warnf(template, args...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Errorw(msg, keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func DPanicw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.DPanicw(msg, keysAndValues...)
}
