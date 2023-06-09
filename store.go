package limiterstore

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type rateLimiter struct {
	limiter   *rate.Limiter
	lastSeen  time.Time
	threshold time.Duration
}

type LimiterStore struct {
	limiterMap      sync.Map
	cleanupInterval time.Duration
	updateInterval  time.Duration
}

func (l *rateLimiter) update(newLimit rate.Limit, newBurst int) {
	if l.limiter.Limit() != newLimit {
		l.limiter.SetLimit(newLimit)
	}
	if l.limiter.Burst() != newBurst {
		l.limiter.SetBurst(newBurst)
	}
}

func NewLimiterStore(cleanupInterval, updateInterval time.Duration) *LimiterStore {
	store := &LimiterStore{
		cleanupInterval: cleanupInterval,
		updateInterval:  updateInterval,
	}
	go store.cleanupUnusedLimiters()
	return store
}

func (s *LimiterStore) LoadAndUpdate(key string, rateLimit rate.Limit, burst int) *rate.Limiter {
	value, loaded := s.limiterMap.LoadOrStore(key, &rateLimiter{
		limiter:   rate.NewLimiter(rateLimit, burst),
		lastSeen:  time.Now(),
		threshold: s.cleanupInterval,
	})

	lim := value.(*rateLimiter)

	// 检查最后一次更新时间是否超过指定的时间间隔，如果超过则执行更新操作
	if loaded && time.Since(lim.lastSeen) > s.updateInterval {
		lim.update(rateLimit, burst)
		lim.lastSeen = time.Now()
	}

	return lim.limiter
}

func (s *LimiterStore) cleanupUnusedLimiters() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.limiterMap.Range(func(key, value interface{}) bool {
			lim := value.(*rateLimiter)
			if time.Since(lim.lastSeen) > lim.threshold {
				s.limiterMap.Delete(key)
			}
			return true
		})
	}
}
