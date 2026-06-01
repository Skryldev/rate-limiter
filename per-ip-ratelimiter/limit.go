package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Middleware func(http.Handler) http.Handler
type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}
type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func perClientRateLimiter() Middleware {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			defer mu.Unlock()

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			if _, found := clients[ip]; !found {
				clients[ip] = &Client{
					limiter: rate.NewLimiter(rate.Every(5*time.Second), 5),
				}
			}
			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				message := Message{
					Status: "Request limit exceeded",
					Body:   "limit exceeded",
				}
				w.Header().Set("X-RateLimit-Limit", "5")
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "5")
				w.WriteHeader(http.StatusTooManyRequests)
				err := json.NewEncoder(w).Encode(&message)
				if err != nil {
					log.Printf("failed to encode error response: %v", err)
					return
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
