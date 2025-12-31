package main

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Limiter struct {
	rate    int
	burst   int
	clients map[string]*visitor
	mu      sync.Mutex
}

type visitor struct {
	limiter  *rate.Limiter
	lastseen time.Time
}

func NewLimiter(rate, burst int) *Limiter {
	return &Limiter{
		rate:    rate,
		burst:   burst,
		clients: make(map[string]*visitor),
	}
}

func (l *Limiter) CleanupClients() {
	for {
		time.Sleep(1 * time.Minute)

		for ip, v := range l.clients {
			if time.Since(v.lastseen) > 3*time.Minute {
				delete(l.clients, ip)
			}
		}
	}
}

func (l *Limiter) GetClient(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, found := l.clients[ip]

	if !found {
		limiter := rate.NewLimiter(rate.Limit(l.rate), l.burst)
		v := &visitor{limiter: limiter, lastseen: time.Now()}
		l.clients[ip] = v
		return limiter
	}
	l.clients[ip].lastseen = time.Now()
	return v.limiter
}

func (app *application) limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return
		}
		limiter := app.limiter.GetClient(ip)
		if limiter.Allow() == false {
			app.rateLimitExceededResponse(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
