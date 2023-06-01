package limiterstore

import (
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestLimiterStore(t *testing.T) {
	cleanupInterval := 5 * time.Second
	updateInterval := 2 * time.Second

	store := NewLimiterStore(cleanupInterval, updateInterval)

	// 创建一个速率限制器并获取
	rateLimit := rate.Limit(100)
	burst := 10
	key := "test-key"
	limiter := store.LoadAndUpdate(key, rateLimit, burst)

	// 验证速率限制器的初始状态
	if limiter.Limit() != rateLimit {
		t.Errorf("Expected limit: %v, got: %v", rateLimit, limiter.Limit())
	}
	if limiter.Burst() != burst {
		t.Errorf("Expected burst: %v, got: %v", burst, limiter.Burst())
	}

	// 等待一段时间超过更新时间间隔，然后再次获取速率限制器
	time.Sleep(2 * updateInterval)
	limiter = store.LoadAndUpdate(key, rateLimit, burst)

	// 验证速率限制器是否已更新
	if limiter.Limit() != rateLimit {
		t.Errorf("Expected limit: %v, got: %v", rateLimit, limiter.Limit())
	}
	if limiter.Burst() != burst {
		t.Errorf("Expected burst: %v, got: %v", burst, limiter.Burst())
	}

	// 等待一段时间超过清理时间间隔，然后再次获取速率限制器
	time.Sleep(2 * cleanupInterval)
	limiter = store.LoadAndUpdate(key, rateLimit, burst)

	// 验证速率限制器是否仍然存在
	if limiter == nil {
		t.Errorf("Expected limiter to exist, got nil")
	}
}
