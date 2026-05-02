package pokecache

import (
	"testing"
	"time"
)

func TestCache_AddAndGet(t *testing.T) {
	cache := NewCache(5 * time.Second)
	defer cache.Close()

	// Test adding and retrieving a value
	testKey := "test-key"
	testValue := []byte("test-value")
	cache.Add(testKey, testValue)

	retrievedValue, found := cache.Get(testKey)
	if !found {
		t.Error("Expected to find the cached value")
	}

	if string(retrievedValue) != string(testValue) {
		t.Errorf("Expected '%s', got '%s'", testValue, retrievedValue)
	}
}

func TestCache_GetNonExistent(t *testing.T) {
	cache := NewCache(5 * time.Second)
	defer cache.Close()

	// Test getting a non-existent key
	_, found := cache.Get("non-existent-key")
	if found {
		t.Error("Expected not to find the cached value")
	}
}

func TestCache_Reaping(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)
	defer cache.Close()

	// Add a value that should expire quickly
	testKey := "expired-key"
	testValue := []byte("expired-value")
	cache.Add(testKey, testValue)

	// Verify it exists initially
	if _, found := cache.Get(testKey); !found {
		t.Error("Expected to find the cached value immediately after adding")
	}

	// Wait for the reaper to run (slightly longer than the interval)
	time.Sleep(150 * time.Millisecond)

	// Verify it's been removed
	if _, found := cache.Get(testKey); found {
		t.Error("Expected the cached value to be reaped after expiration")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(5 * time.Second)
	defer cache.Close()

	// Test concurrent access to the cache
	done := make(chan bool, 100)

	for i := 0; i < 100; i++ {
		go func(i int) {
			key := "concurrent-key"
			value := []byte("value")
			cache.Add(key, value)
			_, _ = cache.Get(key)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}
}
