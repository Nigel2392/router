package middleware

import (
	"encoding/base64"
	"hash/fnv"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Nigel2392/router/v3"
	"github.com/Nigel2392/router/v3/request"
	"golang.org/x/time/rate"
)

// Rate limit types.
type RateLimitType int

// Rate limit types, used to determine what to rate limit with.
const (
	RateLimitIP RateLimitType = iota
	RateLimitIP_Proxy
	RateLimitCookie
)

// Default rate limit options if none are provided.
var defaultConf = &RateLimitOptions{
	RequestsPerSecond: 10,
	BurstMultiplier:   3,
	CleanExpiry:       5 * time.Minute,
	CleanInt:          1 * time.Minute,
	Type:              RateLimitIP,
}

// RateLimitOptions is a struct that holds the options for the rate limiter.
type RateLimitOptions struct {
	CookieName        string
	Type              RateLimitType
	RequestsPerSecond int
	BurstMultiplier   int
	CleanExpiry       time.Duration
	CleanInt          time.Duration
	LimitHandler      func(r *request.Request)
}

// Rate returns the rate limit.
func (r *RateLimitOptions) rate() rate.Limit {
	return rate.Limit(r.RequestsPerSecond)
}

// Burst returns the burst limit.
func (r *RateLimitOptions) burst() int {
	return r.RequestsPerSecond * r.BurstMultiplier
}

// Rate Limiter Middleware
func RateLimiterMiddleware(conf *RateLimitOptions) func(next router.Handler) router.Handler {
	// Use default config if none is provided
	if conf == nil {
		conf = defaultConf
	}

	// Start goroutine to go through and clean up old visitors
	go cleanupVisitors(conf.CleanInt, conf.CleanExpiry)

	// Return the middleware function
	switch conf.Type {
	case RateLimitIP, RateLimitIP_Proxy:
		goto rateLimitByIP
	case RateLimitCookie:
		goto rateLimitByCookie
	default:
		goto rateLimitByIP
	}

	// Rate limit by cookie
rateLimitByCookie:
	return middlewareFunc(func(next router.Handler, r *request.Request) {
		// Get the cookie from the request
		var cookie, err = r.Request.Cookie(conf.CookieName)
		if err != nil {
			// Create a new cookie if it doesn't exist
			cookie = &http.Cookie{
				Name: conf.CookieName,
				// Generate a unique ID for the cookie
				Value:    generateUniqueID(),
				HttpOnly: true,
				Path:     "/",
			}
		}
		cookie.Expires = time.Now().Add(conf.CleanExpiry)
		r.SetCookies(cookie)
		makeChoice(cookie.Value, conf, next, r)
	})

	// Rate limit by IP
rateLimitByIP:
	// Return the middleware function
	return middlewareFunc(func(next router.Handler, r *request.Request) {
		// Get the IP from the request
		var ip string
		var err error
		// If we are rate limiting behind a proxy, get the IP from the headers
		if conf.Type == RateLimitIP_Proxy {
			if ip = r.GetHeader("X-Forwarded-For"); ip != "" {
				r.Request.RemoteAddr = ip
			} else if ip = r.GetHeader("X-Real-IP"); ip != "" {
				r.Request.RemoteAddr = ip
			}
		}
		// Split the IP from the port
		if ip, _, err = net.SplitHostPort(r.Request.RemoteAddr); err != nil {
			if DEFAULT_LOGGER != nil {
				DEFAULT_LOGGER.Error(FormatMessage(r, "ERROR", "Error getting IP: %s", err.Error()))
			}
			r.Error(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		// Choose to allow or disallow the request
		makeChoice(ip, conf, next, r)
	})
}

// Visitor struct, holds the limiter and the last time the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Map of visitors, and a mutex to lock it.
var visitors = make(map[string]*visitor)
var mu = &sync.Mutex{}

// Get a visitor from the map, or create a new one if it doesn't exist.
func getVisitor(id string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	v, ok := visitors[id]
	if !ok {
		limiter := rate.NewLimiter(r, b)
		visitors[id] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// Go through and clean up old visitors.
func cleanupVisitors(clean_interval, clean_expiry time.Duration) {
	var t = time.NewTicker(clean_interval)
	for {
		<-t.C
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > clean_expiry {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// Make a choice, either allow the request or return a 429.
func makeChoice(id string, conf *RateLimitOptions, next router.Handler, r *request.Request) {
	limiter := getVisitor(id, conf.rate(), conf.burst())
	if !limiter.Allow() {
		if conf.LimitHandler != nil {
			conf.LimitHandler(r)
			return
		}
		r.Error(http.StatusTooManyRequests, "Too Many Requests")
		return
	}

	next.ServeHTTP(r)
}

// Shorthand for creating a middleware function.
func middlewareFunc(f func(next router.Handler, r *request.Request)) func(next router.Handler) router.Handler {
	// Return the middleware function
	return func(next router.Handler) router.Handler {
		return router.HandleFunc(func(r *request.Request) {
			f(next, r)
		})
	}

}

var seed = rand.NewSource(time.Now().UnixNano())

// Generate a unique ID.
func generateUniqueID() string {
	var b = make([]byte, 16)
	for i := range b {
		b[i] = byte(seed.Int63())
	}
	return base64.URLEncoding.EncodeToString(fnv.New32a().Sum(b))
}
