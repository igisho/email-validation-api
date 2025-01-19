package main

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var (
    failedLogins = make(map[string]int)
    bannedIPs    = make(map[string]time.Time)
    mu           sync.Mutex
)

// isWhitelisted checks if the IP is whitelisted.
func isWhitelisted(ip string) bool {
    if ip == "127.0.0.1" || ip == "::1" {
        return true
    }
    for _, whitelistedIP := range config.Fail2BanAllow {
        if ip == whitelistedIP {
            return true
        }
    }
    return false
}

// fail2banMiddleware implements fail2ban-like functionality.
func fail2banMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            logError("Error parsing IP address:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        if isWhitelisted(ip) {
            next.ServeHTTP(w, r)
            return
        }

        mu.Lock()
        if banTime, found := bannedIPs[ip]; found && time.Now().Before(banTime) {
            mu.Unlock()
            logError("IP banned:", ip)
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        mu.Unlock()

        next.ServeHTTP(w, r)
    })
}

// trackFailedLogin tracks failed login attempts and bans IP addresses after a threshold.
func trackFailedLogin(r *http.Request) {
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        logError("Error parsing IP address:", err)
        return
    }

    mu.Lock()
    defer mu.Unlock()

    failedLogins[ip]++
    if failedLogins[ip] >= 5 {
        bannedIPs[ip] = time.Now().Add(config.Fail2BanDuration * time.Second)
        logError("IP banned:", ip)
        delete(failedLogins, ip)
    }
}
