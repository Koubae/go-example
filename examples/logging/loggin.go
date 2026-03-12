package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// https://github.com/uber-go/zap
// https://betterstack.com/community/guides/logging/go/zap/

func main() {
	// logger, _ := zap.NewProduction()
	// logger := zap.NewExample()
	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	url := "https://example.com"

	sugar.Infow(
		"failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	// EXAMPLE 2
	logger2 := createLogger()
	defer logger2.Sync()

	logger2.Info("Hello from Zap!")

	// EXAMPLE 3
	logger3 := createLoggerSimple()
	defer logger3.Sync()

	logger3.Info("Hello from Zap!")
	logger3.Warn("Hello from Zap!")
	logger3.Error("Hello from Zap!")

	// EXAMPLE 4
	logger4 := createLoggerFull()
	defer logger4.Sync()

	logger4.Info("Hello from Zap!")
	logger4.Warn("Hello from Zap!")
	logger4.Error("Hello from Zap (I expect stack trace on this one)!")
	logger4.Info("ALL GOOD :=)")

}

func createLoggerSimple() *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	// developmentCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	developmentCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	// developmentCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	// developmentCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	developmentCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		// enc.AppendString(t.UTC().Format(time.RFC3339))                   // always UTC
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z0700")) // always UTC == ISO8601TimeEncoder
	}

	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
	)
	return zap.New(core)
}

func createLoggerFull() *zap.Logger {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // << colors on the level
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encCfg), // must be CONSOLE (not JSON)
		zapcore.Lock(os.Stdout),
		zap.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	return logger
}

func createLogger() *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     7, // days
		},
	)

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core)
}

// func NewCustom() *zap.Logger {
// 	// encoderCfg := zapcore.EncoderConfig{
// 	// 	MessageKey:     "msg",
// 	// 	LevelKey:       "level",
// 	// 	NameKey:        "logger",
// 	// 	EncodeLevel:    zapcore.LowercaseLevelEncoder,
// 	// 	EncodeTime:     zapcore.ISO8601TimeEncoder,
// 	// 	EncodeDuration: zapcore.StringDurationEncoder,
// 	// }
// 	encoderCfg := zap.NewDevelopmentEncoderConfig()
//
// 	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zap.InfoLevel)
// 	return New(core).WithOptions(options...)
// }
//

// / ---------------------------------
//
//	GO Logger setup V2
//
// / ---------------------------------
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

var LogLevels = [4]LogLevel{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError}

func CreateLogger(logLevel LogLevel, filePath string) (*zap.Logger, func(loggerTmp *zap.Logger, logger *zap.Logger)) {
	stdout := zapcore.Lock(zapcore.AddSync(os.Stdout))

	var fileWriter zapcore.WriteSyncer
	if filePath == "" {
		fileWriter = stdout
	} else {
		fileWriter = zapcore.AddSync(createFileLoggerWriter(filePath))
	}

	zapLevel := zapcore.InfoLevel
	switch strings.ToLower(string(logLevel)) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		log.Fatal("Logger level invalid, must be one of: DEBUG, INFO, WARN, or ERROR")
	}
	level := zap.NewAtomicLevelAt(zapLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = utcTimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	developmentCfg.EncodeTime = utcTimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, fileWriter, level),
	)

	return zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel)), logCloser
}

func utcTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
}

func createFileLoggerWriter(filePath string) *lumberjack.Logger {
	logger := zap.L()
	dir := filepath.Dir(filePath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			logger.Panic(fmt.Sprintf("failed to create log directory %q, error: %v", dir, err))
		}
	}

	return &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	}
}

func logCloser(loggerTmp *zap.Logger, logger *zap.Logger) {
	err := logger.Sync()
	if err == nil {
		return
	}

	/*  Ignore some errors related to closing the logger.
	@bug:
		- https://github.com/uber-go/zap/issues/772
		- https://github.com/uber-go/zap/issues/328
	 Also, this seems to work:
	!errors.Is(err, syscall.EINVAL)
	*/
	if _, ok := errors.AsType[*fs.PathError](err); !ok {
		loggerTmp.Error(
			"Error while shutting down logger",
			zap.String("type", fmt.Sprintf("%T", err)),
			zap.Error(err),
		)
	}
}
