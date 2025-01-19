package main

import (
	"sync"
	"time"
)

// Cache structures for MX records and SMTP checks.
type Cache struct {
    MX     map[string][]string
    SMTP   map[string]bool
    Times  map[string]time.Time
    mu     sync.Mutex
    ticker *time.Ticker
}

var cache = Cache{
    MX:     make(map[string][]string),
    SMTP:   make(map[string]bool),
    Times:  make(map[string]time.Time),
    ticker: time.NewTicker(1 * time.Hour),
}

// InitCache initializes the cache and starts the garbage collector.
func InitCache() {
    go func() {
        for range cache.ticker.C {
            cache.GarbageCollect()
        }
    }()
}

// GetMX retrieves MX records from the cache if they haven't expired.
func (c *Cache) GetMX(domain string) ([]string, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if records, found := c.MX[domain]; found {
        if time.Since(c.Times[domain]) < config.CacheMaxTimeMX*time.Second {
            return records, true
        }
    }
    return nil, false
}

// SetMX stores MX records in the cache.
func (c *Cache) SetMX(domain string, records []string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.MX[domain] = records
    c.Times[domain] = time.Now()
}

// GetSMTP retrieves SMTP check results from the cache if they haven't expired.
func (c *Cache) GetSMTP(key string) (bool, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if valid, found := c.SMTP[key]; found {
        if time.Since(c.Times[key]) < config.CacheMaxTimeSMTP*time.Second {
            return valid, true
        }
    }
    return false, false
}

// SetSMTP stores SMTP check results in the cache.
func (c *Cache) SetSMTP(key string, valid bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.SMTP[key] = valid
    c.Times[key] = time.Now()
}

// GarbageCollect removes expired entries from the cache.
func (c *Cache) GarbageCollect() {
    c.mu.Lock()
    defer c.mu.Unlock()

    for key, timestamp := range c.Times {
        if time.Since(timestamp) > config.CacheMaxTimeMX*time.Second && time.Since(timestamp) > config.CacheMaxTimeSMTP*time.Second {
            delete(c.MX, key)
            delete(c.SMTP, key)
            delete(c.Times, key)
        }
    }
}
