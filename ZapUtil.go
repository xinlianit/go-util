// zap日志框架-使用篇 https://zhuanlan.zhihu.com/p/141321801
// zap日志框架-源码篇 https://zhuanlan.zhihu.com/p/141793571
// zap日志框架-性能篇 https://zhuanlan.zhihu.com/p/142196437
package util

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	zapUtilInstance *Zap
	zapUtilOnce     sync.Once
)

func ZapUtil() *Zap {
	zapUtilOnce.Do(func() {
		zapUtilInstance = new(Zap)
	})
	return zapUtilInstance
}

// 切割类型
type RotateType int

const (
	// 按大小切割
	RotateTypeSize RotateType = 1
	// 按日期切割
	RotateTypeDate RotateType = 2
)

// 日志配置
type ZapUtilConfig struct {
	// 日志文件
	LogFile string
	// 错误日志
	ErrorLogFile string
	// 最低记录日志级别
	LowestLevel zapcore.Level
	// 是否记录行号
	RecordLineNumber bool
	// 日志基础字段
	BaseFields []zap.Field
	// 日志格式; 文本格式: zapcore.NewConsoleEncoder、JSON格式: zapcore.NewJSONEncoder
	LogFormatter zapcore.Encoder
	// 是否开启日志切割
	RotateEnable bool
	// 日志切割类型
	RotateType RotateType
	// 大小切割配置
	RotateSizeConfig ZapUtilRotateSizeConfig
	// 日期切割配置
	RotateDateConfig ZapUtilRotateDateConfig
}

// 日志大小切割配置
type ZapUtilRotateSizeConfig struct {
	MaxSize    int  // 在进行切割之前，日志文件的最大大小（以MB为单位)
	MaxBackups int  // 保留旧文件的最大个数
	MaxAge     int  // 保留旧文件的最大天数
	Compress   bool // 是否压缩/归档旧文件
}

// 日志日期切割配置
type ZapUtilRotateDateConfig struct {
	// 切割后缀
	Extend string
	// 日志切割配置
	Options []rotatelogs.Option
}

// Zap 日志记录器
type Zap struct {
	logger *zap.Logger
	config ZapUtilConfig
}

// 初始化
func (u *Zap) Init(config ZapUtilConfig) *zap.Logger {
	// 配置初始化
	u.config = config

	//自定义日志级别：自定义Info级别
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.WarnLevel && level >= u.config.LowestLevel
	})

	//自定义日志级别：自定义Warn级别
	warnLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.WarnLevel && level >= u.config.LowestLevel
	})

	// 日志切割
	//var writeHook io.Writer
	var infoLog, errLog io.Writer
	if !config.RotateEnable {
		// 日志写入
		//writeHook = u.hookDefault(config.LogFile)
		infoLog = u.hookDefault(config.LogFile)
		if config.ErrorLogFile != "" && config.LogFile != config.ErrorLogFile {
			errLog = u.hookDefault(config.ErrorLogFile)
		}
	} else {
		// 切割类型
		switch config.RotateType {
		// 按大小切割
		case RotateTypeSize:
			//writeHook = u.hookSizeRotate(config.LogFile, config.RotateSizeConfig)
			infoLog = u.hookSizeRotate(config.LogFile, config.RotateSizeConfig)
			if config.ErrorLogFile != "" && config.LogFile != config.ErrorLogFile {
				errLog = u.hookSizeRotate(config.ErrorLogFile, config.RotateSizeConfig)
			}
			break
		// 按日期切割
		case RotateTypeDate:
			//writeHook = u.hookDateRotate(config.LogFile, config.RotateDateConfig)
			infoLog = u.hookDateRotate(config.LogFile, config.RotateDateConfig)
			if config.ErrorLogFile != "" && config.LogFile != config.ErrorLogFile {
				errLog = u.hookDateRotate(config.ErrorLogFile, config.RotateDateConfig)
			}
			break
		default:
			//writeHook = u.hookDefault(config.LogFile)
			infoLog = u.hookDefault(config.LogFile)
			if config.ErrorLogFile != "" && config.LogFile != config.ErrorLogFile {
				errLog = u.hookDefault(config.ErrorLogFile)
			}
			break
		}
	}

	// 写入文件
	//write := zapcore.AddSync(writeHook)

	// 创建 zapcore
	//zapCore := zapcore.NewCore(u.config.LogFormatter, write, u.config.LowestLevel)

	// 多个输出
	var cores []zapcore.Core
	// info及以下日志级别记录
	cores = append(cores, zapcore.NewCore(u.config.LogFormatter, zapcore.AddSync(infoLog), infoLevel))
	if config.ErrorLogFile != "" && config.LogFile != config.ErrorLogFile {
		// warn及以上日志级别记录
		cores = append(cores, zapcore.NewCore(u.config.LogFormatter, zapcore.AddSync(errLog), warnLevel))
	}

	// 同时将日志输出到控制台
	cores = append(cores, zapcore.NewCore(u.config.LogFormatter, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), u.config.LowestLevel))

	zapCore := zapcore.NewTee(cores...)

	var options []zap.Option
	// 堆栈跟踪
	options = append(options, zap.AddStacktrace(zap.WarnLevel))

	if u.config.RecordLineNumber {
		// 文件及行号
		options = append(options, zap.AddCaller())
	}

	// 日志基础字段
	if u.config.BaseFields != nil {
		options = append(options, zap.Fields(u.config.BaseFields...))
	}

	// 构造日志
	u.logger = zap.New(zapCore, options...)

	return u.logger
}

// 默认日志输出
func (u *Zap) hookDefault(fileName string) io.Writer {
	var write io.Writer

	// 默认输出到标准输出
	if fileName == "" {
		write = os.Stdout
	} else {
		// 检测目录
		path, _ := os.Getwd()
		dir := filepath.Dir(filepath.Join(path, fileName))
		_, err := os.Stat(dir)
		// 目录不存在
		if err != nil && os.IsNotExist(err) {
			// 创建目录
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				panic(err)
			}
		}

		// 日志文件
		write, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}

	return write
}

// 文件大小切割 Hook
func (u *Zap) hookSizeRotate(fileName string, config ZapUtilRotateSizeConfig) io.Writer {
	// 日志写入
	writeHook := &lumberjack.Logger{
		Filename:   fileName,          // 日志文件的位置
		MaxSize:    config.MaxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位)
		MaxBackups: config.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     config.MaxAge,     // 保留旧文件的最大天数
		Compress:   config.Compress,   // 是否压缩/归档旧文件
	}

	return writeHook
}

// 日期切割 Hook
func (u *Zap) hookDateRotate(fileName string, config ZapUtilRotateDateConfig) io.Writer {
	extend := u.config.RotateDateConfig.Extend
	options := u.config.RotateDateConfig.Options
	options = append(options, rotatelogs.WithLinkName(fileName)) // 生成软链，指向最新日志文件
	writeHook, err := rotatelogs.New(fileName+extend, options...)

	if err != nil {
		// 日志切割创建失败 todo 记录日志
		//u.Logger().Errorf("Hook Error: %v", err.Error())
		panic(err)
	}

	return writeHook
}

// 默认编码器配置
func (u *Zap) NewDefaultEncoderConfig() zapcore.EncoderConfig {
	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "file"
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stacktrace"

	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		// 时间格式化
		encoder.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	// 级别名称大写
	//encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return encoderConfig
}

// 默认大小切割配置
func (u *Zap) NewDefaultRotateSizeConfig() ZapUtilRotateSizeConfig {
	return ZapUtilRotateSizeConfig{
		MaxSize:    10,
		MaxBackups: 100,
		MaxAge:     30,
		Compress:   true,
	}
}

// 默认日期切割配置
func (u *Zap) NewDefaultRotateDateConfig() ZapUtilRotateDateConfig {
	return ZapUtilRotateDateConfig{
		Extend: ".%Y%m%d",
		Options: []rotatelogs.Option{
			rotatelogs.WithLinkName(""),                 // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(30 * 24 * time.Hour),  // 设置最大保存时间(30天)
			rotatelogs.WithRotationTime(24 * time.Hour), // 设置日志切割时间间隔(1天)
		},
	}
}

// 默认配置
func (u *Zap) NewDefaultConfig() ZapUtilConfig {
	// 编码器配置
	encoderConfig := u.NewDefaultEncoderConfig()

	// 编码器
	//encoder := zapcore.NewJSONEncoder(encoderConfig) // JSON格式
	encoder := zapcore.NewConsoleEncoder(encoderConfig) // 文本格式

	return ZapUtilConfig{
		LogFile:          "",
		ErrorLogFile:     "",
		LowestLevel:      zapcore.DebugLevel,
		RecordLineNumber: true,
		BaseFields:       nil,
		LogFormatter:     encoder,
		RotateEnable:     false,
		RotateType:       0,
		RotateSizeConfig: u.NewDefaultRotateSizeConfig(),
		RotateDateConfig: u.NewDefaultRotateDateConfig(),
	}
}

// 日志记录器
func (u *Zap) Logger() *zap.Logger {
	return u.logger
}
