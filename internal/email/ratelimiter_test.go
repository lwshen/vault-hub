package email

import (
	"testing"
	"time"
)

func TestMemoryRateLimiterAllow(t *testing.T) {
	limiter := newMemoryRateLimiter(50*time.Millisecond, 2)
	key := "test@example.com"

	if !limiter.Allow(key) {
		t.Fatal("expected first attempt to be allowed")
	}
	if !limiter.Allow(key) {
		t.Fatal("expected second attempt to be allowed within burst")
	}
	if limiter.Allow(key) {
		t.Fatal("expected third attempt to be denied due to rate limit")
	}
}

func TestMemoryRateLimiterResetsAfterWindow(t *testing.T) {
	limiter := newMemoryRateLimiter(20*time.Millisecond, 1)
	key := "reset@example.com"

	if !limiter.Allow(key) {
		t.Fatal("expected first attempt to be allowed")
	}
	if limiter.Allow(key) {
		t.Fatal("expected second attempt to be denied before window expires")
	}

	time.Sleep(25 * time.Millisecond)

	if !limiter.Allow(key) {
		t.Fatal("expected attempt after window to be allowed")
	}
}

func TestMemoryRateLimiterEmptyKey(t *testing.T) {
	limiter := newMemoryRateLimiter(10*time.Millisecond, 1)
	if !limiter.Allow("") {
		t.Fatal("expected Allow to succeed for empty key")
	}
}
