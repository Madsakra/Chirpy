package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"example.com/m/v2/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	plat := os.Getenv("PLATFORM")
	log.Println("DB_URL is:", dbURL)
	dbs, err := sql.Open("postgres", dbURL)

	if err != nil {
		panic(err)
	}
	defer dbs.Close()
	dbQueries := database.New(dbs)

	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       plat,
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/chirps", apiCfg.handleVerification)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirp)

	mux.HandleFunc("POST /api/users", apiCfg.CreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.Login)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
