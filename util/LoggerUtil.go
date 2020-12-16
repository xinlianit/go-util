package util

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	logUtilInstance *Logger
	logUtilOnce     sync.Once
)

func LoggerUtil() *Logger {
	logUtilOnce.Do(func() {
		logUtilInstance = new(Logger)
	})
	return logUtilInstance
}

// 日志配置
type LoggerUtilConfig struct {
	// 日志文件
	LogFile string
	// 最大记录日志级别
	LowestLevel logrus.Level
	// 日志格式
	LogFormatter logrus.Formatter
	// 是否开启日志切割
	RotateEnable bool
	// 切割后缀
	RotateExtend string
	// 日志切割配置
	RotateOptions []rotatelogs.Option
}

// 日志工具
type Logger struct {
	logger *logrus.Logger
	config LoggerUtilConfig
}

// 初始化
func (u *Logger) Init(config LoggerUtilConfig) *logrus.Logger {
	// 默认配置
	if config.LogFormatter == nil {
		// 默认配置
		config.LogFormatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		}
	}

	u.config = config
	u.logger = u.newFileLogger(u.config.LogFile)

	return u.logger
}

// 创建文件日志记录器
func (u *Logger) newFileLogger(fileName string) *logrus.Logger {
	// 日志记录器实例
	logger := logrus.New()

	// 设置最低loglevel
	logger.SetLevel(u.config.LowestLevel)

	if !u.config.RotateEnable {
		// 设置日志输出output（无锁模式）
		logger.Out = u.createLogOut(fileName)

		// 设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
		logger.SetFormatter(u.config.LogFormatter)
	} else {
		logWrite, err := rotatelogs.New(fileName+u.config.RotateExtend, u.config.RotateOptions...)

		if err != nil {
			// 日志切割创建失败
			u.Logger().Errorf("Hook Error: %v", err.Error())
			return nil
		}

		writeMap := lfshook.WriterMap{
			logrus.DebugLevel: logWrite,
			logrus.InfoLevel:  logWrite,
			logrus.WarnLevel:  logWrite,
			logrus.ErrorLevel: logWrite,
			logrus.FatalLevel: logWrite,
			logrus.PanicLevel: logWrite,
		}

		hook := lfshook.NewHook(writeMap, u.config.LogFormatter)

		logger.AddHook(hook)
	}

	return logger
}

// 创建日志输出
func (u *Logger) createLogOut(fileName string) *os.File {
	// todo 递归创建日志目录
	// 日志输出
	if fileName != "" {
		// 日志文件
		logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		return logFile
	}

	// 默认输出控制台
	return os.Stdout
}

// 默认记录器
func (u *Logger) Logger() *logrus.Logger {
	return u.logger
}

// 记录日志到文件
func (u *Logger) LoggerToFile(fileName string) *logrus.Logger {
	return u.newFileLogger(fileName)
}
