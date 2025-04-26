package auth

import (
	"time"
	"sync"
	"net/http"
)

type rateLimiter struct {
	visits map[string]int
	mu     sync.Mutex
}

var limiter = rateLimiter{visits: make(map[string]int)}

// RateLimitMiddleware limits the number of requests from a single IP.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter.mu.Lock()
		count := limiter.visits[ip]
		if count >= 5 {
			limiter.mu.Unlock()
			http.Error(w, "Too many requests. Try again later.", http.StatusTooManyRequests)
			return
		}
		limiter.visits[ip] = count + 1
		limiter.mu.Unlock()

		// Reset the count after 15 minutes
		go func(ip string) {
			time.Sleep(15 * time.Minute)
			limiter.mu.Lock()
			delete(limiter.visits, ip)
			limiter.mu.Unlock()
		}(ip)

		next.ServeHTTP(w, r)
	})
}