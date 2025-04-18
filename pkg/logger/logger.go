package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// LogLevel 定义日志级别类型
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

var logger = logrus.New()

// 自定义格式化器，用于简化日志输出
type simpleFormatter struct {
	showCaller bool
}

func (f *simpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := strings.ToUpper(entry.Level.String())
	padding := "     "
	level = level + padding[len(level):]

	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	callerInfo := ""
	if f.showCaller && entry.Caller != nil {
		// 获取相对路径，简化调用信息
		fileName := filepath.Base(entry.Caller.File)
		callerInfo = fmt.Sprintf(" [%s:%d]", fileName, entry.Caller.Line)
	}

	// 如果消息以 [Manager] 等前缀开头，不再添加额外前缀
	message := entry.Message

	return []byte(fmt.Sprintf("%s | %s%s | %s\n", level, timestamp, callerInfo, message)), nil
}

func init() {
	// 设置自定义格式化器
	logger.SetFormatter(&simpleFormatter{showCaller: false})

	// 设置输出
	logger.SetOutput(os.Stdout)

	// 设置日志级别
	logger.SetLevel(logrus.InfoLevel)

	// 不显示调用位置，由自定义格式化器处理
	logger.SetReportCaller(false)

	// 尝试从环境变量获取日志级别
	logLevelEnv := os.Getenv("LOG_LEVEL")
	if logLevelEnv != "" {
		switch strings.ToLower(logLevelEnv) {
		case "debug":
			SetLevel(logrus.DebugLevel)
		case "info":
			SetLevel(logrus.InfoLevel)
		case "warn", "warning":
			SetLevel(logrus.WarnLevel)
		case "error":
			SetLevel(logrus.ErrorLevel)
		case "fatal":
			SetLevel(logrus.FatalLevel)
		}
	}
}

// SetOutput 设置日志输出位置
func SetOutput(out io.Writer) {
	logger.SetOutput(out)
}

// SetLevel 设置日志级别
func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

// Debug 输出调试级别日志
func Debug(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Info 输出信息级别日志
func Info(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn 输出警告级别日志
func Warn(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Error 输出错误级别日志
func Error(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal 输出致命错误日志并退出程序
func Fatal(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// WithField 添加单个字段
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

// WithFields 添加多个字段
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}
