package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// LogLevel 日志级别
type LogLevel int

const (
	// DEBUG 调试级别
	DEBUG LogLevel = iota
	// INFO 信息级别
	INFO
	// WARN 警告级别
	WARN
	// ERROR 错误级别
	ERROR
	// FATAL 致命错误级别
	FATAL
)

var (
	// 日志实例
	logger = logrus.New()
)

// init 初始化日志配置
func init() {
	// 设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		DisableColors:          false,
		ForceColors:            true,
		DisableQuote:           false,
		QuoteEmptyFields:       true,
		PadLevelText:           true,
		FieldMap:               nil,
		CallerPrettyfier:       nil,
	})

	// 设置输出
	logger.SetOutput(os.Stdout)

	// 设置日志级别
	logger.SetLevel(logrus.InfoLevel)

	// 显示调用位置
	logger.SetReportCaller(true)
}
