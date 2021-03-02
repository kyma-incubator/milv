package pkg

import "time"

type retry struct {
	d time.Duration
}

func NewRetry(d time.Duration) *retry {
	return &retry{d: d}
}

func (r *retry) Limit() {
	time.Sleep(r.d)
}
