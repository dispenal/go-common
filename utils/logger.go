package common_utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func init() {
	var err error

	var encoder zapcore.Encoder

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.NameKey = "[SERVICE]"
	encoderConfig.TimeKey = "[TIME]"
	encoderConfig.LevelKey = "[LEVEL]"
	encoderConfig.CallerKey = "[LINE]"
	encoderConfig.MessageKey = "[MESSAGE]"
	encoderConfig.StacktraceKey = os.Getenv("SERVICE_NAME")
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	logWriter := zapcore.AddSync(os.Stdout)

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	encoderConfig.ConsoleSeparator = " | "

	encoder = zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(loggerLevelMap["debug"]))

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	if err != nil {
		logger.Fatal("Error when init logger: " + err.Error())
	}
}

func LogInfo(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func LogFatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

func LogPanic(message string, fields ...zap.Field) {
	logger.Panic(message, fields...)
}

func LogDebug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func LogError(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}
