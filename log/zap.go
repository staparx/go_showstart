package log

import (
	"fmt"
	"os"
	"path"
	"time"

	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var level zapcore.Level

func InitLogger() {
	Logger = initZap()
	return
}

var director = "log/log_file"

func initZap() (logger *zap.Logger) {

	if ok, _ := PathExists("log/log_file"); !ok { // 判断是否有Director文件夹
		fmt.Printf("创建日志文件 %v \n", director)
		_ = os.Mkdir(director, os.ModePerm)
	}

	level = zap.InfoLevel
	logger = zap.New(getEncoderCore())
	//logger = logger.WithOptions(zap.AddCaller())

	return logger
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	//普通日志和错误日志分开存储
	//新加一个debug日志，用于调试
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= level
	})
	errLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.InfoLevel
	})

	infoWriter, err := getWriteSyncer("info")
	if err != nil {
		fmt.Printf("Get infoWriter Syncer Failed err:%v", err.Error())
		return
	}
	errWriter, err := getWriteSyncer("error")
	if err != nil {
		fmt.Printf("Get errWriter Syncer Failed err:%v", err.Error())
		return
	}
	debugWriter, err := getWriteSyncer("debug")
	if err != nil {
		fmt.Printf("Get debugWriter Syncer Failed err:%v", err.Error())
		return
	}

	return zapcore.NewTee(
		zapcore.NewCore(getEncoder(), infoWriter, infoLevel),
		zapcore.NewCore(getEncoder(), errWriter, errLevel),
		zapcore.NewCore(getEncoder(), debugWriter, debugLevel),
	)
}

func getWriteSyncer(fileName string) (zapcore.WriteSyncer, error) {
	zaprotatelogs.ForceNewFile()
	fileWriter, err := zaprotatelogs.New(
		path.Join(director, fileName+".log"),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	// 如果为debug日志，不输出到控制台
	if fileName == "debug" {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(fileWriter)), err
	} else {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}

}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	//if config.Cfg.Zap.Format == "json" {
	//	return zapcore.NewJSONEncoder(getEncoderConfig())
	//}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (zapConfig zapcore.EncoderConfig) {
	zapConfig = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "log",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	return zapConfig
}

// CustomTimeEncoder 自定义日志输出时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000"))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
