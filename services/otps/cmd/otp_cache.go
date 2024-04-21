package main

import (
	"sync"
	"time"
)

type otpCache struct {
	mu   *sync.Mutex
	otps map[string]time.Time
}

func (c *otpCache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for otp, createdAt := range c.otps {
		if time.Since(createdAt) > 30*time.Second {
			delete(c.otps, otp)
		}
	}
}

func (c *otpCache) get(otp string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.otps[otp]
	return ok
}

func (c *otpCache) add(otp string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.otps[otp] = time.Now()
}

func newCache() *otpCache {
	return &otpCache{
		mu:   &sync.Mutex{},
		otps: map[string]time.Time{},
	}
}
