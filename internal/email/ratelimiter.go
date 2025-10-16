package email

import (
	"errors"
	"sync"
	"time"

	"github.com/lwshen/vault-hub/internal/config"
)

var ErrRateLimited = errors.New("email send rate limited")

type RateLimiter interface {
	Allow(key string) bool
}

type noopRateLimiter struct{}

func (noopRateLimiter) Allow(string) bool { return true }

type memoryRateLimiter struct {
	mu     sync.Mutex
	window time.Duration
	max    int
	hits   map[string][]time.Time
}

func newMemoryRateLimiter(window time.Duration, max int) *memoryRateLimiter {
	return &memoryRateLimiter{
		window: window,
		max:    max,
		hits:   make(map[string][]time.Time),
	}
}

func (m *memoryRateLimiter) Allow(key string) bool {
	if key == "" || m.max <= 0 {
		return true
	}

	now := time.Now()
	cutoff := now.Add(-m.window)

	m.mu.Lock()
	defer m.mu.Unlock()

	timestamps := m.hits[key]
	var filtered []time.Time
	if len(timestamps) > 0 {
		filtered = timestamps[:0]
		for _, ts := range timestamps {
			if ts.After(cutoff) {
				filtered = append(filtered, ts)
			}
		}
	}

	if len(filtered) >= m.max {
		m.hits[key] = filtered
		return false
	}

	filtered = append(filtered, now)
	m.hits[key] = filtered
	return true
}

var defaultRateLimiter RateLimiter = noopRateLimiter{}

func DefaultRateLimiter() RateLimiter {
	return defaultRateLimiter
}

func SetDefaultRateLimiter(l RateLimiter) {
	if l == nil {
		defaultRateLimiter = noopRateLimiter{}
		return
	}
	defaultRateLimiter = l
}

func init() {
	if !config.EmailEnabled || !config.EmailRateLimitEnabled {
		return
	}
	if config.EmailRateLimitWindow <= 0 || config.EmailRateLimitBurst <= 0 {
		return
	}
	defaultRateLimiter = newMemoryRateLimiter(config.EmailRateLimitWindow, config.EmailRateLimitBurst)
}
