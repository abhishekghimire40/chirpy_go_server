package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/abhishekghimire40/chirpy_go_server/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	const rootFilePath = "."
	const port = ":8080"
	// main router
	router := chi.NewRouter()

	// creating a new database connection
	rootDir, _ := os.Getwd()

	filepath := filepath.Join(rootDir, "database.json")
	db, err := database.NewDB(filepath)
	if err != nil {
		log.Fatal(err)
	}
	// handling different urls
	apiCfg := &apiConfig{
		fileServerHits: 0,
	}
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootFilePath))))
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)
	router.Handle("/assets/logo", http.FileServer(http.Dir("./assets")))

	// apiRouter
	apiRouter := chi.NewRouter()

	apiRouter.Get("/healthz", handleReadiness)
	apiRouter.Post("/chirps", ValidateChirp(db))
	apiRouter.Get("/chirps", GetAllChirps(db))
	apiRouter.Get("/chirps/{chirpID}", GetSingleChirp(db))
	router.Mount("/api", apiRouter)

	// adminRouter
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.HandleMetrics)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	// different features of url
	srv := &http.Server{
		Addr:    port,
		Handler: corsMux,
	}
	log.Printf("Serving on port: %s", port)
	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Fatal("Error Startng server: ", err)
	}
}
