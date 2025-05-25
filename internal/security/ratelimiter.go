package sec

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/utkukrl/safegate/internal/core"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(rateStr string, burst int) *RateLimiter {
	r := parseRate(rateStr)
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    burst,
	}
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func RateLimitMiddleware(rl *RateLimiter) core.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := core.ClientIP(r).String()
			limiter := rl.getVisitor(ip)
			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func parseRate(s string) rate.Limit {
	if strings.HasSuffix(s, "r/s") {
		val := strings.TrimSuffix(s, "r/s")
		n, err := strconv.Atoi(val)
		if err == nil && n > 0 {
			return rate.Limit(n)
		}
	}
	return rate.Every(time.Second)
}
