package pkg

import "time"

type limiter struct {
	d time.Duration
}

func NewLimiter(d time.Duration) *limiter {
	return &limiter{d: d}
}

func (r *limiter) Limit() {
	time.Sleep(r.d)
}
