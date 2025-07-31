package ratelimiter

import (
	"golang.org/x/time/rate"
	"sync"
)

type TokenBucket struct {
	mu      sync.Mutex
	limiter map[string]*rate.Limiter
}

func NewTokenBucket() *TokenBucket {
	return &TokenBucket{
		limiter: make(map[string]*rate.Limiter),
	}
}
func (tb *TokenBucket) GetLimiter(key string, r rate.Limit, b int) *rate.Limiter {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if limiter, exists := tb.limiter[key]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(r, b)
	tb.limiter[key] = limiter
	return limiter
}
