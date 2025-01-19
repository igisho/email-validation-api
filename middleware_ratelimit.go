package main

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// rateLimitMiddleware implements rate limiting.
func rateLimitMiddleware(next http.Handler) http.Handler {
    limiter := rate.NewLimiter(rate.Every(time.Minute/time.Duration(config.RequestsPerMinute)), config.RequestsPerMinute)
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            logError("Rate limit exceeded")
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
