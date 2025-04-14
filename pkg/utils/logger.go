package utils

import (
	"log"
	"os"
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
	// 当前日志级别，默认为INFO
	currentLevel = INFO
	// 日志前缀
	debugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger  = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger  = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(os.Stdout, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)
)

// SetLogLevel 设置日志级别
func SetLogLevel(level LogLevel) {
	currentLevel = level
}

// Debug 输出调试级别日志
func Debug(format string, v ...interface{}) {
	if currentLevel <= DEBUG {
		debugLogger.Printf(format, v...)
	}
}

// Info 输出信息级别日志
func Info(format string, v ...interface{}) {
	if currentLevel <= INFO {
		infoLogger.Printf(format, v...)
	}
}

// Warn 输出警告级别日志
func Warn(format string, v ...interface{}) {
	if currentLevel <= WARN {
		warnLogger.Printf(format, v...)
	}
}

// Error 输出错误级别日志
func Error(format string, v ...interface{}) {
	if currentLevel <= ERROR {
		errorLogger.Printf(format, v...)
	}
}

// Fatal 输出致命错误级别日志并退出程序
func Fatal(format string, v ...interface{}) {
	fatalLogger.Printf(format, v...)
	os.Exit(1)
}
