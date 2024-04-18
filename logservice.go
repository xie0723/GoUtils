package GoUtils

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func genEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func genInfoWrite() zapcore.WriteSyncer {
	infoLumberIO := &lumberjack.Logger{
		Filename:   "./logs/info.log",
		MaxSize:    10, // megabytes
		MaxBackups: 100,
		MaxAge:     7, // days
		Compress:   false,
	}
	return zapcore.AddSync(infoLumberIO)
}

func genErrorWrite() zapcore.WriteSyncer {
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   "./logs/error.log",
		MaxSize:    10, // megabytes
		MaxBackups: 100,
		MaxAge:     7, // days
		Compress:   false,
	}
	return zapcore.AddSync(lumberWriteSyncer)
}

func InitLogger() *zap.SugaredLogger {
	encoder := genEncoder()

	info := genInfoWrite()
	err := genErrorWrite()

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(info, zapcore.AddSync(os.Stdout)), zap.DebugLevel),
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(err, zapcore.AddSync(os.Stdout)), zap.ErrorLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
	return logger.Sugar()
}
