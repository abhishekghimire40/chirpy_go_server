package main

import (
	"fmt"
	"net/http"
	"sync"
)

// apiConfig struct for counting number of times a certain path is called in api
type apiConfig struct {
	fileServerHits int
	hitsMutex      sync.Mutex
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.hitsMutex.Lock()
		defer c.hitsMutex.Unlock()
		c.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf("Hits: %d", c.fileServerHits)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
