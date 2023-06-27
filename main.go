package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	const rootFilePath = "."
	const port = ":8080"
	//
	r := chi.NewRouter()
	// handling different urls
	apiCfg := &apiConfig{
		fileServerHits: 0,
	}
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootFilePath))))
	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)
	r.Handle("/assets/logo", http.FileServer(http.Dir("./assets")))
	r.Get("/healthz", handleReadiness)
	r.Get("/metrics", apiCfg.HandleMetrics)

	corsMux := middlewareCors(r)

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
