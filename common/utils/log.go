package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// 创建日志器
func NewLogger(serviceName string) *logrus.Entry {
	logger := logrus.New()

	// 设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志级别
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// 返回带有服务名称的日志条目
	return logger.WithField("service", serviceName)
}
