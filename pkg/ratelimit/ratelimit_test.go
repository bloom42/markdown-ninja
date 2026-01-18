package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimit_Basic(t *testing.T) {
	l := New()
	defer l.Stop()

	actor := []byte("user123")
	action := "test-action"
	bucket := time.Minute
	allowed := uint64(3)

	// First 3 should be allowed
	for i := 0; i < 3; i++ {
		if !l.RateLimit(action, actor, bucket, allowed) {
			t.Errorf("Request %d should have been allowed", i+1)
		}
	}

	// 4th should be rate limited
	if l.RateLimit(action, actor, bucket, allowed) {
		t.Error("Request 4 should have been rate limited")
	}

	// 5th should also be rate limited
	if l.RateLimit(action, actor, bucket, allowed) {
		t.Error("Request 5 should have been rate limited")
	}
}

func TestRateLimit_DifferentActors(t *testing.T) {
	l := New()
	defer l.Stop()

	actor1 := []byte("user1")
	actor2 := []byte("user2")
	action := "test-action"
	bucket := time.Minute
	allowed := uint64(2)

	// Actor 1 uses both requests
	l.RateLimit(action, actor1, bucket, allowed)
	l.RateLimit(action, actor1, bucket, allowed)

	// Actor 1 should be limited
	if l.RateLimit(action, actor1, bucket, allowed) {
		t.Error("Actor1 should be rate limited")
	}

	// Actor 2 should still be allowed
	if !l.RateLimit(action, actor2, bucket, allowed) {
		t.Error("Actor2 should be allowed")
	}
}

func TestRateLimit_DifferentActions(t *testing.T) {
	l := New()
	defer l.Stop()

	actor := []byte("user123")
	action1 := "action1"
	action2 := "action2"
	bucket := time.Minute
	allowed := uint64(1)

	// Use up action1 limit
	l.RateLimit(action1, actor, bucket, allowed)

	// Action1 should be limited
	if l.RateLimit(action1, actor, bucket, allowed) {
		t.Error("Action1 should be rate limited")
	}

	// Action2 should still be allowed
	if !l.RateLimit(action2, actor, bucket, allowed) {
		t.Error("Action2 should be allowed")
	}
}

func TestRateLimit_BucketReset(t *testing.T) {
	l := New()
	defer l.Stop()

	actor := []byte("user123")
	action := "test-action"
	bucket := 100 * time.Millisecond
	allowed := uint64(2)

	// Use up the limit
	l.RateLimit(action, actor, bucket, allowed)
	l.RateLimit(action, actor, bucket, allowed)

	if l.RateLimit(action, actor, bucket, allowed) {
		t.Error("Should be rate limited")
	}

	// Wait for next bucket
	time.Sleep(150 * time.Millisecond)

	// Should be allowed again
	if !l.RateLimit(action, actor, bucket, allowed) {
		t.Error("Should be allowed after bucket reset")
	}
}

func TestCount(t *testing.T) {
	l := New()
	defer l.Stop()

	actor := []byte("user123")
	action := "test-action"
	bucket := time.Minute

	if c := l.Count(action, actor, bucket); c != 0 {
		t.Errorf("Expected count 0, got %d", c)
	}

	l.RateLimit(action, actor, bucket, 10)
	l.RateLimit(action, actor, bucket, 10)

	if c := l.Count(action, actor, bucket); c != 2 {
		t.Errorf("Expected count 2, got %d", c)
	}
}

func TestRemaining(t *testing.T) {
	l := New()
	defer l.Stop()

	actor := []byte("user123")
	action := "test-action"
	bucket := time.Minute
	allowed := uint64(5)

	if r := l.Remaining(action, actor, bucket, allowed); r != 5 {
		t.Errorf("Expected 5 remaining, got %d", r)
	}

	l.RateLimit(action, actor, bucket, allowed)
	l.RateLimit(action, actor, bucket, allowed)

	if r := l.Remaining(action, actor, bucket, allowed); r != 3 {
		t.Errorf("Expected 3 remaining, got %d", r)
	}
}

func TestRateLimit_Concurrent(t *testing.T) {
	l := New()
	defer l.Stop()

	actor := []byte("user123")
	action := "test-action"
	bucket := time.Minute
	allowed := uint64(100)

	var wg sync.WaitGroup
	allowedCount := 0
	var mu sync.Mutex

	// Run 200 concurrent requests
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.RateLimit(action, actor, bucket, allowed) {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Exactly 100 should have been allowed
	if allowedCount != 100 {
		t.Errorf("Expected 100 allowed, got %d", allowedCount)
	}
}

func TestRateLimit_BinaryActor(t *testing.T) {
	l := New()
	defer l.Stop()

	// Test with binary data containing special characters
	actor := []byte{0x00, 0x01, 0xFF, ':', '\n'}
	action := "test-action"
	bucket := time.Minute
	allowed := uint64(1)

	if !l.RateLimit(action, actor, bucket, allowed) {
		t.Error("First request should be allowed")
	}
	if l.RateLimit(action, actor, bucket, allowed) {
		t.Error("Second request should be rate limited")
	}
}
