package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Request struct {
	UserID  string `json:"user_id"`
	Payload string `json:"payload"`
}

type StatsResponse struct {
	Stats map[string]int `json:"stats"`
}

type RateLimiter struct {
	mu     sync.Mutex
	users  map[string][]time.Time
	totals map[string]int
}

func (rl *RateLimiter) Allow(userID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times := rl.users[userID]

	var valid []time.Time
	for _, t := range times {
		if now.Sub(t) < time.Minute {
			valid = append(valid, t)
		}
	}

	if len(valid) >= 5 {
		return false
	}

	valid = append(valid, now)
	rl.users[userID] = valid
	rl.totals[userID]++
	return true
}

func (rl *RateLimiter) GetStats() map[string]int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	stats := make(map[string]int)
	for k, v := range rl.totals {
		stats[k] = v
	}
	return stats
}

var limiter = &RateLimiter{
	users:  make(map[string][]time.Time),
	totals: make(map[string]int),
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if !limiter.Allow(req.UserID) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Process payload - for now, just accept
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request accepted"))
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := limiter.GetStats()
	resp := StatsResponse{Stats: stats}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/request", handleRequest)
	http.HandleFunc("/stats", handleStats)
	http.ListenAndServe(":8080", nil)
}
