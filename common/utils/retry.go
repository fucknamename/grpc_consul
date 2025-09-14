package utils

import (
	"math"
	"math/rand"
	"time"
)

// 重试配置
type RetryConfig struct {
	MaxRetries     int           // 最大重试次数
	InitialBackoff time.Duration // 初始退避时间
	MaxBackoff     time.Duration // 最大退避时间
	Multiplier     float64       // 退避倍数
}

// 默认重试配置
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:     10,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
	}
}

// 带退避的重试函数
func RetryWithBackoff(operation func() error, config *RetryConfig) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error
	for attempt := 0; attempt < config.MaxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}
		lastErr = err

		if attempt < config.MaxRetries-1 {
			backoff := calculateBackoff(attempt, config)
			time.Sleep(backoff)
		}
	}

	return lastErr
}

// 计算退避时间
func calculateBackoff(attempt int, config *RetryConfig) time.Duration {
	// 指数退避
	backoff := float64(config.InitialBackoff) * math.Pow(config.Multiplier, float64(attempt))

	// 添加抖动以避免雷鸣群效应
	jitter := rand.Float64() * 0.3 * backoff
	backoff = backoff + jitter

	// 确保不超过最大退避时间
	if backoff > float64(config.MaxBackoff) {
		backoff = float64(config.MaxBackoff)
	}

	return time.Duration(backoff)
}

// 永久重试直到成功
func RetryForever(operation func() error, backoffConfig *RetryConfig) {
	if backoffConfig == nil {
		backoffConfig = DefaultRetryConfig()
	}

	attempt := 0
	for {
		err := operation()
		if err == nil {
			return
		}

		backoff := calculateBackoff(attempt, backoffConfig)
		time.Sleep(backoff)

		// 限制attempt增长以保持最大退避时间
		if attempt < 10 {
			attempt++
		}
	}
}
