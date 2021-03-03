package pkg

import "time"

type waiter struct {
	d time.Duration
}

func NewWaiter(d time.Duration) *waiter {
	return &waiter{d: d}
}

func (l *waiter) Wait() {
	time.Sleep(l.d)
}
