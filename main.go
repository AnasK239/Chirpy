package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/AnasK239/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries *database.Queries
	platform string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM") 
	
	db , err := sql.Open("postgres" , dbURL)
	if err != nil{
		log.Fatalf("Error opening database conneciton")
	}
	
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		dbQueries: database.New(db),
		platform: platform,
	}

	serveMux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	FileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	FileHandler = apiCfg.middleWareMetricsIncr(FileHandler)

	serveMux.Handle("/app/", FileHandler)

	serveMux.HandleFunc("GET /api/healthz", HealthHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)

	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	serveMux.HandleFunc("POST /api/users" , apiCfg.handlerCreateUser)
	serveMux.HandleFunc("POST /api/chirps" , apiCfg.handleCreateChirp)
	
	serveMux.HandleFunc("GET /api/chirps" , apiCfg.handleGetAllChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}" , apiCfg.handleGetChirp)

	serveMux.HandleFunc("POST /api/login" , apiCfg.handleLogin)


	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
