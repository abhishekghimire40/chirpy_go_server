package main

import (
	"fmt"
	"net/http"
	"sync"
)

// apiConfig struct for counting number of times a certain path is called in api
type apiConfig struct {
	fileServerHits int
	jwtSecret      string
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
	response := fmt.Sprintf(`
	<html>
	<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
	</body>
	</html>
	`, c.fileServerHits)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
