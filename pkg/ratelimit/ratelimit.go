// Package ratelimit provides a simple fixed-window rate limiter.
package ratelimit

import (
	"encoding/binary"
	"sync"
	"time"

	"github.com/skerkour/stdx-go/crypto/blake3"
)

// Limiter tracks request counts within time buckets.
type Limiter struct {
	mutex   sync.Mutex
	buckets map[[32]byte]*bucket
	stop    chan struct{}
}

type bucket struct {
	count   uint64
	expires time.Time
}

// New creates a new rate limiter with automatic cleanup of expired buckets.
func New() *Limiter {
	limiter := &Limiter{
		mutex:   sync.Mutex{},
		buckets: make(map[[32]byte]*bucket),
		stop:    make(chan struct{}),
	}
	go limiter.cleanupLoop()
	return limiter
}

// RateLimit checks if an action by an actor is allowed within the rate limit.
// It returns true if the action is allowed, false if rate limited.
//
// Parameters:
//   - action: identifies the type of action being rate limited (e.g., "login", "api-call")
//   - actor: identifies who is performing the action (e.g., user ID, IP address)
//   - timeBucket: the duration of each rate limit window
//   - allowed: maximum number of actions allowed per time bucket
func (limiter *Limiter) RateLimit(action string, actor []byte, timeBucket time.Duration, allowed uint64) bool {
	now := time.Now()
	bucketStart := now.Truncate(timeBucket)
	key := makeKey(action, actor, uint64(bucketStart.UnixNano()), uint64(timeBucket.Nanoseconds()))

	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	b, exists := limiter.buckets[key]
	if !exists {
		limiter.buckets[key] = &bucket{
			count:   1,
			expires: bucketStart.Add(timeBucket * 2), // Keep for one extra period for safety
		}
		return true
	}

	if b.count >= allowed {
		return false
	}

	b.count++
	return true
}

// Count returns the current count for an action/actor in the current time bucket.
// Useful for showing users how many requests they have remaining.
func (limiter *Limiter) Count(action string, actor []byte, timeBucket time.Duration) uint64 {
	now := time.Now()
	bucketStart := now.Truncate(timeBucket)
	key := makeKey(action, actor, uint64(bucketStart.UnixNano()), uint64(timeBucket.Nanoseconds()))

	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	if b, exists := limiter.buckets[key]; exists {
		return b.count
	}
	return 0
}

// Remaining returns how many requests are remaining for an action/actor.
func (limiter *Limiter) Remaining(action string, actor []byte, timeBucket time.Duration, allowed uint64) uint64 {
	count := limiter.Count(action, actor, timeBucket)
	if count >= allowed {
		return 0
	}
	return allowed - count
}

// Stop stops the background cleanup goroutine.
// Call this when the limiter is no longer needed.
func (limiter *Limiter) Stop() {
	close(limiter.stop)
}

func (limiter *Limiter) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			limiter.cleanup()
		case <-limiter.stop:
			return
		}
	}
}

func (limiter *Limiter) cleanup() {
	now := time.Now()

	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	for key, b := range limiter.buckets {
		if now.After(b.expires) {
			delete(limiter.buckets, key)
		}
	}
}

func makeKey(action string, actor []byte, bucketStartNanos uint64, timeBucketNanos uint64) [32]byte {
	var hash [32]byte

	hasher := blake3.New(32, nil)
	hasher.Write([]byte(action))
	hasher.Write(actor)
	binary.Write(hasher, binary.LittleEndian, bucketStartNanos)
	binary.Write(hasher, binary.LittleEndian, timeBucketNanos)

	hasher.Sum(hash[:0])
	return hash
}
