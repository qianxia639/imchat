package utils

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Zap(logpath, loglevel string) *zap.Logger {

	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	case "warn":
		level = zap.WarnLevel
	default:
		level = zap.InfoLevel
	}

	write := zapcore.AddSync(newLumberjack(logpath))

	encodreConfig := customEncoderConfig()

	consoleEncoder := zapcore.NewConsoleEncoder(encodreConfig)
	fileEncoder := zapcore.NewJSONEncoder(encodreConfig)

	cores := make([]zapcore.Core, 0)

	// 开发环境，同时在控制台输出
	if level == zap.DebugLevel {
		core := zapcore.NewCore(consoleEncoder, os.Stdout, level)
		cores = append(cores, core)
	}

	cores = append(cores, zapcore.NewCore(fileEncoder, write, level))
	core := zapcore.NewTee(cores...)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// // 开启文件及行号
	development := zap.Development()
	return zap.New(core, caller, development)
}

// 自定义打印格式
func customEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		NameKey:     "logger",
		CallerKey:   "linenum",
		FunctionKey: zapcore.OmitKey,
		MessageKey:  "msg",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
			pae.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z0700"))
		}, // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 短路径编码
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// 日志切割
func newLumberjack(logpath string) io.Writer {
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s%s.log", logpath, time.Now().UTC().Format("2006-01-02")), // 日志文件路径，路径为空时会使用 os.TempDir()
		MaxSize:    100,                                                                     // 文件最大大小，默认100M
		MaxBackups: 30,                                                                      // 日志文件保存最大数量，默认全部保存
		MaxAge:     7,                                                                       // 保存的最大天数，默认不限
		Compress:   true,                                                                    // 是否压缩，默认不压缩
	}
}
