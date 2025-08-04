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
