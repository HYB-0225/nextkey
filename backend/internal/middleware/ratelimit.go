package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// 定期清理过期记录
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, timestamps := range rl.requests {
		var valid []time.Time
		for _, t := range timestamps {
			if now.Sub(t) < rl.window {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = valid
		}
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// 获取该IP的请求记录
	timestamps, exists := rl.requests[ip]
	if !exists {
		rl.requests[ip] = []time.Time{now}
		return true
	}

	// 过滤掉窗口外的请求
	var validTimestamps []time.Time
	for _, t := range timestamps {
		if now.Sub(t) < rl.window {
			validTimestamps = append(validTimestamps, t)
		}
	}

	// 检查是否超过限制
	if len(validTimestamps) >= rl.limit {
		return false
	}

	// 添加新请求
	validTimestamps = append(validTimestamps, now)
	rl.requests[ip] = validTimestamps

	return true
}

var loginLimiter = newRateLimiter(5, 1*time.Minute)

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := newRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.allow(ip) {
			utils.Error(c, 429, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoginRateLimitMiddleware 登录专用速率限制中间件
func LoginRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !loginLimiter.allow(ip) {
			utils.Error(c, 429, "登录尝试过于频繁，请1分钟后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
