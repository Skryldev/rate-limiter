package main

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func rateLimiter(next http.HandlerFunc) http.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(5*time.Second), 5)

	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := Message{
				Status: "Request limit exceeded",
				Body:   "limit exceeded",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(w).Encode(&message)
			if err != nil {
				return
			}
			return
		}

		next(w, r)
	}
}
