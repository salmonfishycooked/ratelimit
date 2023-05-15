package ratelimit

import (
	"sync"
	"time"
)

const (
	DEFAULT_QUANTUM = 1
)

type Bucket struct {
	// rate of token production, nums of quantum token is generated per unit(rate) time
	rate time.Duration

	quantum int64

	// capacity of bucket
	cap int64

	// the number of tokens remaining in the bucket
	tokens int64

	// the last time tokens were put in
	latestTime time.Time

	mu sync.Mutex
}

// Allow used to determine whether a token is currently available
func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.adjustAvailableToken()
	if b.tokens > 0 {
		b.tokens--
		return true
	} else {
		return false
	}
}

// adjustAvailableToken used to adjust the number of tokens
func (b *Bucket) adjustAvailableToken() {
	now := time.Now()
	nums := int64(now.Sub(b.latestTime)/b.rate) * b.quantum
	if nums == 0 {
		return
	}

	b.latestTime = now
	b.tokens += nums
	if b.tokens > b.cap {
		b.tokens = b.cap
	}
	b.latestTime = time.Now()
}

func NewBucket(rate time.Duration, cap int64) *Bucket {
	return NewBucketWithQuantum(rate, cap, DEFAULT_QUANTUM)
}

func NewBucketWithQuantum(rate time.Duration, cap int64, quantum int64) *Bucket {
	return &Bucket{
		rate:       rate,
		quantum:    quantum,
		cap:        cap,
		tokens:     cap,
		latestTime: time.Now(),
	}
}
