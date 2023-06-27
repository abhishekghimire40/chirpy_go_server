package main

import (
	"log"
	"net/http"
)

func main() {
	const rootFilePath = "."
	const port = ":8080"
	//
	mux := http.NewServeMux()
	// handling different urls
	apiCfg := &apiConfig{
		fileServerHits: 0,
	}
	metricsHandler := http.FileServer(http.Dir(rootFilePath))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", metricsHandler)))
	mux.Handle("/assets/logo", http.FileServer(http.Dir("./assets")))
	mux.HandleFunc("/healthz", handleReadiness)
	mux.HandleFunc("/metrics", apiCfg.HandleMetrics)

	corsMux := middlewareCors(mux)

	// different features of url
	srv := &http.Server{
		Addr:    port,
		Handler: corsMux,
	}
	log.Printf("Serving on port: %s", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error Startng server: ", err)
	}

}
