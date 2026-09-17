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

// http.FileServer(http.Dir(".")) is a handler that serves files from the exact URL resource path given
// http.NewServeMux creates an Interception point for requests
// calling .Handle on a serveMux requires giving the request path to be handled
// and the handler that handles it
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
	
	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
