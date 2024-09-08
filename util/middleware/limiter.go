package middleware

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type UserLimiter struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var (
	limiters = make(map[string]*UserLimiter)
	mutex    sync.Mutex
)

const cleanupInterval = time.Minute * 5

func GetLimiter(userID string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()

	if limiter, exists := limiters[userID]; exists {
		limiter.LastSeen = time.Now()
		return limiter.Limiter
	}

	limiter := rate.NewLimiter(rate.Every(time.Minute*5), 5)
	limiters[userID] = &UserLimiter{
		Limiter:  limiter,
		LastSeen: time.Now(),
	}
	return limiter
}

func CleanupOldLimiters() {
	for {
		time.Sleep(cleanupInterval)

		mutex.Lock()
		for userID, limiter := range limiters {
			if time.Since(limiter.LastSeen) > cleanupInterval {
				delete(limiters, userID)
			}
		}
		mutex.Unlock()
	}
}
