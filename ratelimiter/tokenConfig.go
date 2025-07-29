package ratelimiter

import "time"

type TokenBucket struct {
	tokens    int
	max       int
	refilRate time.Duration
}
